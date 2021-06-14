package network

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/automatico/jato/internal/utils"
	"github.com/automatico/jato/pkg/constant"
	"github.com/automatico/jato/pkg/data"
	"golang.org/x/crypto/ssh"
)

type SSHParams struct {
	Port                int  `json:"port"`
	InsecureConnection  bool `json:"insecureConnection"`
	InsecureCyphers     bool `json:"insecureCyphers"`
	InsecureKeyExchange bool `json:"insecureKeyExchange"`
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

func SSHClientConfig(username string, password string, insecureConnection bool, insecureCyphers bool, InsecureKeyExchange bool) *ssh.ClientConfig {
	c := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	if insecureConnection {
		c.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	if insecureCyphers {
		c.Config.Ciphers = append(c.Config.Ciphers, constant.InsecureSSHCyphers...)
	}
	if InsecureKeyExchange {
		c.KeyExchanges = append(c.KeyExchanges, constant.InsecureSSHKeyAlgorithms...)
	}
	return c
}

func InitSSHParams(s *SSHParams) {
	if s.Port == 0 {
		s.Port = constant.SSHPort
	}
}

func ConnectWithSSH(host string, port int, clientConfig *ssh.ClientConfig) SSHConn {

	sshConn := SSHConn{}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), clientConfig)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		fmt.Println(err)
	}

	stdOut, err := session.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	stdIn, err := session.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = session.RequestPty("xterm", 0, 200, modes)
	if err != nil {
		session.Close()
		fmt.Println(err)
	}

	err = session.Shell()
	if err != nil {
		session.Close()
		fmt.Println(err)
	}

	sshConn.Session = session
	sshConn.StdIn = stdIn
	sshConn.StdOut = stdOut

	return sshConn

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
				// fmt.Println(tmp)
			}
			br <- tmp
		}(stdOut, expect, buffRead)

		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeout) * time.Second):
			fmt.Printf("Waiting for '%s' took longer than timeout: %d\n", expect, timeout)
		}
	}(stdOut, expect)

	return <-ch
}

// RunWithSSH is the entrypoint to run commands
// against a device.
func RunWithSSH(sd SSHDevice, commands []string, ch chan data.Result, wg *sync.WaitGroup) {
	err := sd.ConnectWithSSH()
	if err != nil {
		fmt.Println(err)
	}
	defer sd.DisconnectSSH()
	defer wg.Done()

	result := sd.SendCommandsWithSSH(commands)

	ch <- result
}
