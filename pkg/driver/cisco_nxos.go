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
	CiscoNXOSUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\\-_@()/:]{1,63}>\s$`)
	CiscoNXOSSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}#\s$`)
	CiscoNXOSConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#\s$`)
)

// CiscoNXOSDevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoNXOSDevice struct {
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

func (cd *CiscoNXOSDevice) ConnectWithSSH() error {

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

	cd.SendCommandWithSSH("terminal length 0")
	cd.SendCommandWithSSH("terminal width 511")

	return nil
}

func (cd CiscoNXOSDevice) DisconnectSSH() error {
	return cd.SSHConn.Session.Close()
}

func (cd CiscoNXOSDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandWithSSH(cd.SSHConn, command, cd.SuperUserPromptRE, 5)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (cd CiscoNXOSDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandsWithSSH(cd.SSHConn, commands, cd.SuperUserPromptRE, 5)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewCiscoNXOSDevice takes a NetDevice and initializes
// a CiscoNXOSDevice.
func NewCiscoNXOSDevice(nd NetDevice) CiscoNXOSDevice {
	cd := CiscoNXOSDevice{}
	cd.IP = nd.IP
	cd.Name = nd.Name
	cd.Vendor = nd.Vendor
	cd.Platform = nd.Platform
	cd.Connector = nd.Connector
	cd.Credentials = nd.Credentials
	cd.SSHParams = nd.SSHParams

	// Prompts
	cd.UserPromptRE = CiscoNXOSUserPromptRE
	cd.SuperUserPromptRE = CiscoNXOSSuperUserPromptRE
	cd.ConfigPromtRE = CiscoNXOSConfigPromptRE

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
