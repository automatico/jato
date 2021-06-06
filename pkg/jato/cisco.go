package jato

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/reiver/go-telnet"
	"golang.org/x/crypto/ssh"
)

var (
	CiscoUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\\-_@()/:]{1,63}>$`)
	CiscoSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\\-_@()/:]{1,63}#$`)
	CiscoConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}\([a-z0-9.\-@/:\+]{0,32}\)#$`)
	CiscoDisablePaging     string         = "terminal length 0"
)

type CiscoIOSDevice struct {
	IP                string `json:"ip"`
	Name              string `json:"name"`
	Vendor            string `json:"vendor"`
	Platform          string `json:"platform"`
	Connector         string `json:"connector"`
	UserPromptRE      *regexp.Regexp
	SuperUserPromptRE *regexp.Regexp
	ConfigPromtRE     *regexp.Regexp
	DisablePaging     string
	Credentials
	SSHParams
	TelnetParams
	SSHConn
	TelnetConn *telnet.Conn
}

func (cd *CiscoIOSDevice) Init() {
	// Prompts
	if cd.UserPromptRE == nil {
		cd.UserPromptRE = CiscoUserPromptRE
	}
	if cd.SuperUserPromptRE == nil {
		cd.SuperUserPromptRE = CiscoSuperUserPromptRE
	}
	if cd.ConfigPromtRE == nil {
		cd.ConfigPromtRE = CiscoConfigPromptRE
	}

	// Paging
	if cd.DisablePaging == "" {
		cd.DisablePaging = CiscoDisablePaging
	}

	// SSH Params
	if cd.SSHParams.Port == 0 {
		cd.SSHParams.Port = SSHPort
	}
	if !cd.SSHParams.InsecureConnection {
		cd.SSHParams.InsecureConnection = true
	}
	if !cd.SSHParams.InsecureCyphers {
		cd.SSHParams.InsecureCyphers = true
	}

	// Telnet Params
	if cd.TelnetParams.Port == 0 {
		cd.TelnetParams.Port = TelnetPort
	}
}

func (cd *CiscoIOSDevice) ConnectWithTelnet() error {

	conn, err := telnet.DialTo(fmt.Sprintf("%s:%d", cd.IP, cd.TelnetParams.Port))
	if err != nil {
		return err
	}

	_, err = SendCommandWithTelnet(conn, cd.Username, PasswordRE, 1)
	if err != nil {
		fmt.Println(err)
	}
	_, err = SendCommandWithTelnet(conn, cd.Password, cd.SuperUserPromptRE, 1)
	if err != nil {
		fmt.Println(err)
	}
	_, err = SendCommandWithTelnet(conn, cd.DisablePaging, cd.SuperUserPromptRE, 1)
	if err != nil {
		fmt.Println(err)
	}

	cd.TelnetConn = conn

	return nil
}

func (cd CiscoIOSDevice) DisconnectTelnet() error {
	return cd.TelnetConn.Close()
}

func (cd CiscoIOSDevice) SendCommandWithTelnet(cmd string) Result {

	result := Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := SendCommandWithTelnet(cd.TelnetConn, cmd, cd.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (cd CiscoIOSDevice) SendCommandsWithTelnet(commands []string) Result {

	result := Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := SendCommandsWithTelnet(cd.TelnetConn, commands, cd.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

func (cd *CiscoIOSDevice) ConnectWithSSH() error {

	sshConn := SSHConn{}

	clientConfig := SSHClientConfig(
		cd.Credentials.Username,
		cd.Credentials.Password,
		cd.SSHParams.InsecureConnection,
		cd.SSHParams.InsecureCyphers,
	)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cd.IP, cd.SSHParams.Port), clientConfig)
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

	readSSH(stdOut, cd.SuperUserPromptRE, 2)

	sshConn.Session = session
	sshConn.StdIn = stdIn
	sshConn.StdOut = stdOut

	cd.SSHConn = sshConn

	cd.SendCommandWithSSH(cd.DisablePaging)

	return nil
}

func (cd CiscoIOSDevice) DisconnectSSH() error {
	return cd.SSHConn.Session.Close()
}

func (cd CiscoIOSDevice) SendCommandWithSSH(command string) Result {

	result := Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := SendCommandWithSSH(cd.SSHConn, command, cd.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (cd CiscoIOSDevice) SendCommandsWithSSH(commands []string) Result {

	result := Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := SendCommandsWithSSH(cd.SSHConn, commands, cd.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}
