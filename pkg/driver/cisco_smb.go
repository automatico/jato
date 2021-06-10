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
	CiscoSMBUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\\-_@()/:]{1,63}>$`)
	CiscoSMBSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\\-_@()/:]{1,63}#$`)
	CiscoSMBConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}\([a-z0-9.\-@/:\+]{0,32}\)#$`)
)

// CiscoSMBDevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoSMBDevice struct {
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

func (cd *CiscoSMBDevice) ConnectWithSSH() error {

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

	network.ReadSSH(stdOut, cd.SuperUserPromptRE, 5)

	sshConn.Session = session
	sshConn.StdIn = stdIn
	sshConn.StdOut = stdOut

	cd.SSHConn = sshConn

	cd.SendCommandWithSSH("terminal datadump")
	cd.SendCommandWithSSH("terminal width 512")

	return nil
}

func (cd CiscoSMBDevice) DisconnectSSH() error {
	return cd.SSHConn.Session.Close()
}

func (cd CiscoSMBDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	// Cisco SMB devices are really slow to output to the terminal.
	cmdOut, err := network.SendCommandWithSSH(cd.SSHConn, command, cd.SuperUserPromptRE, 120)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (cd CiscoSMBDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	// Cisco SMB devices are really slow to output to the terminal.
	cmdOut, err := network.SendCommandsWithSSH(cd.SSHConn, commands, cd.SuperUserPromptRE, 120)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewCiscoSMBDevice takes a NetDevice and initializes
// a CiscoSMBDevice.
func NewCiscoSMBDevice(nd NetDevice) CiscoSMBDevice {
	cd := CiscoSMBDevice{}
	cd.IP = nd.IP
	cd.Name = nd.Name
	cd.Vendor = nd.Vendor
	cd.Platform = nd.Platform
	cd.Connector = nd.Connector
	cd.Credentials = nd.Credentials
	cd.SSHParams = nd.SSHParams

	// Prompts
	cd.UserPromptRE = CiscoSMBUserPromptRE
	cd.SuperUserPromptRE = CiscoSMBSuperUserPromptRE
	cd.ConfigPromtRE = CiscoSMBConfigPromptRE

	// SSH Params
	if cd.SSHParams.Port == 0 {
		cd.SSHParams.Port = constant.SSHPort
	}
	if !cd.SSHParams.InsecureConnection {
		cd.SSHParams.InsecureConnection = true
	}
	if !cd.SSHParams.InsecureCyphers {
		cd.SSHParams.InsecureCyphers = true
	}
	if !cd.SSHParams.InsecureKeyExchange {
		cd.SSHParams.InsecureKeyExchange = true
	}

	return cd
}
