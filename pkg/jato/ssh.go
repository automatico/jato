package jato

import (
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/automatico/jato/internal/utils"
	"golang.org/x/crypto/ssh"
)

var (
	SSHPort int = 22
)

type SSHParams struct {
	Port               int
	InsecureConnection bool
	InsecureCyphers    bool
}

type SSHConn struct {
	Session *ssh.Session
	StdIn   io.Writer
	StdOut  io.Reader
}

type SSHDevice interface {
	ConnectWithSSH() error
	SendCommandsWithSSH([]string) Result
	DisconnectSSH() error
}

func SSHClientConfig(username string, password string, insecureConnection bool, insecureCyphers bool) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	if insecureConnection {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	if insecureCyphers {
		config.Config.Ciphers = append(config.Config.Ciphers, "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-cbc", "aes192-cbc", "aes256-cbc", "3des-cbc", "des-cbc")
	}
	return config
}

func SendCommandsWithSSH(conn SSHConn, commands []string, expect *regexp.Regexp, timeout int64) ([]CommandOutput, error) {

	cmdOut := []CommandOutput{}

	for _, cmd := range commands {
		res, err := SendCommandWithSSH(conn, cmd, expect, timeout)
		if err != nil {
			return cmdOut, err
		}
		cmdOut = append(cmdOut, res)
	}

	return cmdOut, nil

}

func SendCommandWithSSH(conn SSHConn, cmd string, expect *regexp.Regexp, timeout int64) (CommandOutput, error) {
	cmdOut := CommandOutput{}

	_, err := writeSSH(conn.StdIn, cmd)
	time.Sleep(time.Millisecond * 3)

	res := readSSH(conn.StdOut, expect, timeout)
	if err != nil {
		return cmdOut, err
	}

	cmdOut.Command = cmd
	cmdOut.CommandU = utils.Underscorer(cmd)
	cmdOut.Output = utils.CleanOutput(res)

	return cmdOut, nil
}

func writeSSH(stdIn io.Writer, cmd string) (int, error) {
	i, err := stdIn.Write([]byte(cmd + "\r"))
	return i, err
}

func readSSH(stdOut io.Reader, expect *regexp.Regexp, timeout int64) string {
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

func RunWithSSH(sd SSHDevice, commands []string, ch chan Result, wg *sync.WaitGroup) {
	err := sd.ConnectWithSSH()
	if err != nil {
		fmt.Println(err)
	}
	defer sd.DisconnectSSH()
	defer wg.Done()

	result := sd.SendCommandsWithSSH(commands)

	ch <- result
}
