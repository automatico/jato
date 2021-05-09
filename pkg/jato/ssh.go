package jato

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/automatico/jato/internal"
	"golang.org/x/crypto/ssh"
)

type SSHConnection struct {
	Session *ssh.Session
	StdIn   io.WriteCloser
	StdOut  io.Reader
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

// Expect like interface
func SSHExpecter(sshConn SSHConnection, cmd string, expect string, timeout int) (string, error) {
	if _, err := writeBuff(cmd, sshConn.StdIn); err != nil {
		return "", err
	}
	result := readBuff(expect, sshConn.StdOut, timeout)

	return result, nil
}

func readBuffForString(expect string, stdOut io.Reader, buffRead chan<- string) {
	buf := make([]byte, 4096)
	n, err := stdOut.Read(buf) //this reads the ssh terminal
	tmp := ""
	if err == nil {
		tmp = string(buf[:n])
	}
	for (err == nil) && (!strings.Contains(tmp, expect)) {
		n, err = stdOut.Read(buf)
		tmp += string(buf[:n])
		// Uncommenting this might help you debug if you are coming into
		// errors with timeouts when correct details entered
		// fmt.Println(tmp)
	}
	buffRead <- tmp
}

func readBuff(expect string, stdOut io.Reader, timeout int) string {
	ch := make(chan string)

	go func(expect string, stdOut io.Reader) {

		buffRead := make(chan string)

		go readBuffForString(expect, stdOut, buffRead)

		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeout) * time.Second):
			fmt.Printf("Waiting for '%s' took longer than timeout: %d\n", expect, timeout)
		}
	}(expect, stdOut)

	return <-ch
}

func writeBuff(cmd string, stdIn io.WriteCloser) (int, error) {
	returnCode, err := stdIn.Write([]byte(cmd + "\r"))
	return returnCode, err
}

func SSHRunner(nd NetDevice, commands []string, ch chan Result, wg *sync.WaitGroup) {

	conn := nd.ConnectWithSSH()
	defer conn.Session.Close()
	defer wg.Done()

	result := Result{}
	cmdOut := []CommandOutput{}

	result.Device = nd.Name
	result.Timestamp = time.Now().Unix()
	for _, command := range commands {
		res, err := SSHExpecter(conn, command, "#", 5)
		if err != nil {
			result.OK = false
			ch <- result
			return
		}
		out := CommandOutput{Command: internal.Underscorer(command), Output: res}
		cmdOut = append(cmdOut, out)
	}
	result.CommandOutputs = cmdOut
	result.OK = true
	ch <- result
}
