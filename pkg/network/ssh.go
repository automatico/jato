package network

import (
	"fmt"
	"io"
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
		config.Config.Ciphers = append(config.Config.Ciphers,
			"aes128-ctr",
			"aes192-ctr",
			"aes256-ctr",
			"aes128-cbc",
			"aes192-cbc",
			"aes256-cbc",
			"3des-cbc",
			"des-cbc")
	}
	if InsecureKeyExchange {
		config.KeyExchanges = append(
			config.KeyExchanges,
			"diffie-hellman-group-exchange-sha256",
			"diffie-hellman-group-exchange-sha1",
			"diffie-hellman-group1-sha1",
			"diffie-hellman-group14-sha1",
		)
	}
	return config
}

func InitSSHParams(s *SSHParams) {
	if s.Port == 0 {
		s.Port = constant.SSHPort
	}
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

func RunWithSSH(sd SSHDevice, commands []string, ch chan data.Result, wg *sync.WaitGroup) {
	fmt.Printf("%+v\n", sd)
	err := sd.ConnectWithSSH()
	if err != nil {
		fmt.Println(err)
	}
	defer sd.DisconnectSSH()
	defer wg.Done()

	result := sd.SendCommandsWithSSH(commands)

	ch <- result
}
