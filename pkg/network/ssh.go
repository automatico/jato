package network

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/automatico/jato/internal/logger"
	"github.com/automatico/jato/internal/utils"
	"github.com/automatico/jato/pkg/constant"
	"github.com/automatico/jato/pkg/data"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type SSHParams struct {
	Port                int    `json:"port"`
	KeyBasedAuth        bool   `json:"keyBasedAuth"`
	PasswordBasedAuth   bool   `json:"passwordBasedAuth"`
	KnownHostsFile      string `json:"knownHostsFile"`
	InsecureConnection  bool   `json:"insecureConnection"`
	InsecureCyphers     bool   `json:"insecureCyphers"`
	InsecureKeyExchange bool   `json:"insecureKeyExchange"`
}

type SSHConn struct {
	Session *ssh.Session
	StdIn   io.Writer
	StdOut  io.Reader
}

type SSHDevice interface {
	ConnectWithSSH() error
	SendCommandsWithSSH([]string) data.Result
	DisconnectSSH() error
}

func SSHClientConfig(c data.Credentials, s SSHParams) *ssh.ClientConfig {
	conf := &ssh.ClientConfig{
		User: c.Username,
	}

	if s.KnownHostsFile == "" {
		s.KnownHostsFile = constant.SSHKnownHostsFile
	}

	// Setup  authentication method
	if c.SSHKeyFile != "" { // prefer connection with an SSH key file
		key, err := ioutil.ReadFile(c.SSHKeyFile)
		if err != nil {
			logger.Fatalf("unable to read private key: %v", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			logger.Fatalf("unable to parse private key: %v", err)
		}
		conf.Auth = append(conf.Auth, ssh.PublicKeys(signer))

	}
	// if no ssh key use password auth
	if c.Password == "" {
		logger.Fatal("an SSH key or password is required.")
	}
	conf.Auth = append(conf.Auth, ssh.Password(c.Password))

	// Setup remote machines host key checking
	if s.InsecureConnection { // NOT RECOMMENDED FOR PRODUCTION
		conf.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		hostKeyCallback, err := knownhosts.New(s.KnownHostsFile)
		if err != nil {
			logger.Fatalf("could not create hostkeycallback function: %v", err)
		}
		conf.HostKeyCallback = hostKeyCallback
	}

	if s.InsecureCyphers {
		conf.Config.Ciphers = append(conf.Config.Ciphers, constant.InsecureSSHCyphers...)
	}

	if s.InsecureKeyExchange {
		conf.KeyExchanges = append(conf.KeyExchanges, constant.InsecureSSHKeyAlgorithms...)
	}

	return conf
}

func InitSSHParams(s *SSHParams) {
	if s.Port == 0 {
		s.Port = constant.SSHPort
	}
	if s.KnownHostsFile == "" {
		s.KnownHostsFile = constant.SSHKnownHostsFile
	}
}

// TODO: createKnownHosts should create the known hosts
// file if it does not exist.
// https://cyruslab.net/2020/10/23/golang-how-to-write-ssh-hostkeycallback/
func createKnownHosts(s string) {
	if _, err := os.Stat(s); os.IsNotExist(err) {
		f, err := os.OpenFile(s, os.O_CREATE, 0600)
		if err != nil {
			logger.Fatal(err)
		}
		f.Close()
	}
}

func ConnectWithSSH(host string, port int, clientConfig *ssh.ClientConfig) (SSHConn, error) {

	sshConn := SSHConn{}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), clientConfig)
	if err != nil {
		return sshConn, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return sshConn, err
	}

	stdOut, err := session.StdoutPipe()
	if err != nil {
		return sshConn, err
	}

	stdIn, err := session.StdinPipe()
	if err != nil {
		return sshConn, err
	}

	err = session.RequestPty("xterm", 0, 200, modes)
	if err != nil {
		session.Close()
		return sshConn, err
	}

	err = session.Shell()
	if err != nil {
		session.Close()
		return sshConn, err
	}

	sshConn.Session = session
	sshConn.StdIn = stdIn
	sshConn.StdOut = stdOut

	return sshConn, nil

}

func SendCommandsWithSSH(conn SSHConn, commands []string, expect *regexp.Regexp, timeout int64) ([]data.CommandOutput, error) {

	cmdOut := []data.CommandOutput{}

	for _, cmd := range commands {
		res, err := SendCommandWithSSH(conn, cmd, expect, timeout)
		if err != nil {
			return cmdOut, err
		}
		cmdOut = append(cmdOut, res)
	}

	return cmdOut, nil

}

func SendCommandWithSSH(conn SSHConn, cmd string, expect *regexp.Regexp, timeout int64) (data.CommandOutput, error) {
	cmdOut := data.CommandOutput{}

	_, err := WriteSSH(conn.StdIn, cmd)
	time.Sleep(time.Millisecond * 3)

	res := ReadSSH(conn.StdOut, expect, timeout)
	if err != nil {
		return cmdOut, err
	}

	cmdOut.Command = cmd
	cmdOut.CommandU = utils.Underscorer(cmd)
	cmdOut.Output = utils.CleanOutput(res)

	return cmdOut, nil
}

func WriteSSH(stdIn io.Writer, cmd string) (int, error) {
	i, err := stdIn.Write([]byte(cmd + "\r"))
	return i, err
}

func ReadSSH(stdOut io.Reader, expect *regexp.Regexp, timeout int64) string {
	ch := make(chan string)

	go func(stdOut io.Reader, expect *regexp.Regexp) {

		buffRead := make(chan string)

		go func(r io.Reader, exp *regexp.Regexp, br chan<- string) {
			buf := make([]byte, 8192)
			n, err := r.Read(buf) //this reads the ssh terminal
			tmp := ""
			if err == nil {
				tmp = string(buf[:n])
			}
			for (err == nil) && !exp.MatchString(tmp) {
				n, err = r.Read(buf)
				tmp += string(buf[:n])
				// Uncommenting this might help you debug if you are coming into
				// errors with timeouts when correct details entered
				// logger.Debug(tmp)
			}
			br <- tmp
		}(stdOut, expect, buffRead)

		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeout) * time.Second):
			logger.Errorf("Waiting for '%s' took longer than timeout: %d", expect, timeout)
		}
	}(stdOut, expect)

	return <-ch
}

// RunWithSSH is the entrypoint to run commands
// against a device.
func RunWithSSH(sd SSHDevice, commands []string, ch chan data.Result, wg *sync.WaitGroup) {
	err := sd.ConnectWithSSH()
	if err != nil {
		logger.Error(err)
	}
	defer sd.DisconnectSSH()
	defer wg.Done()

	result := sd.SendCommandsWithSSH(commands)

	ch <- result
}
