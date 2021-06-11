package driver

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/constant"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
	"golang.org/x/crypto/ssh"
)

var (
	CiscoXRUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}#\s?$`)
	CiscoXRSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}#\s?$`)
	CiscoXRConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#$`)
)

// CiscoIOSXRDevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoIOSXRDevice struct {
	IP                string `json:"ip"`
	Name              string `json:"name"`
	Vendor            string `json:"vendor"`
	Platform          string `json:"platform"`
	Connector         string `json:"connector"`
	UserPromptRE      *regexp.Regexp
	SuperUserPromptRE *regexp.Regexp
	ConfigPromtRE     *regexp.Regexp
	data.Credentials
	network.SSHParams
	network.SSHConn
}

func (cd *CiscoIOSXRDevice) ConnectWithSSH() error {

	sshConn := network.SSHConn{}

	clientConfig := network.SSHClientConfig(
		cd.Credentials.Username,
		cd.Credentials.Password,
		cd.SSHParams.InsecureConnection,
		cd.SSHParams.InsecureCyphers,
		cd.SSHParams.InsecureKeyExchange,
	)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
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

	network.ReadSSH(stdOut, cd.SuperUserPromptRE, 2)

	sshConn.Session = session
	sshConn.StdIn = stdIn
	sshConn.StdOut = stdOut

	cd.SSHConn = sshConn

	cd.SendCommandWithSSH("terminal length 0")
	cd.SendCommandWithSSH("terminal width 0")

	return nil
}

func (cd CiscoIOSXRDevice) DisconnectSSH() error {
	return cd.SSHConn.Session.Close()
}

func (cd CiscoIOSXRDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandWithSSH(cd.SSHConn, command, cd.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (cd CiscoIOSXRDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandsWithSSH(cd.SSHConn, commands, cd.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewCiscoIOSXRDevice takes a NetDevice and initializes
// a CiscoIOSXRDevice.
func NewCiscoIOSXRDevice(nd NetDevice) CiscoIOSXRDevice {
	cd := CiscoIOSXRDevice{}
	cd.IP = nd.IP
	cd.Name = nd.Name
	cd.Vendor = nd.Vendor
	cd.Platform = nd.Platform
	cd.Connector = nd.Connector
	cd.Credentials = nd.Credentials
	cd.SSHParams = nd.SSHParams

	// Prompts
	cd.UserPromptRE = CiscoXRUserPromptRE
	cd.SuperUserPromptRE = CiscoXRSuperUserPromptRE
	cd.ConfigPromtRE = CiscoXRConfigPromptRE

	// SSH Params
	if cd.SSHParams.Port == 0 {
		cd.SSHParams.Port = constant.SSHPort
	}
	if !cd.SSHParams.InsecureConnection {
		cd.SSHParams.InsecureConnection = true
	}
	if !cd.SSHParams.InsecureCyphers {
		cd.SSHParams.InsecureCyphers = false
	}
	if !cd.SSHParams.InsecureKeyExchange {
		cd.SSHParams.InsecureKeyExchange = false
	}
	return cd
}
