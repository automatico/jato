package driver

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
	"golang.org/x/crypto/ssh"
)

var (
	CiscoAireOSUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)^\([a-z0-9.\\-_\s@()/:]{1,63}\)\s>$`)
	CiscoAireOSSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)^\([a-z0-9.\\-_\s@()/:]{1,63}\)\s>$`)
	CiscoAireOSConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)^\([a-z0-9.\\-_\s@()/:]{1,63}\)\sconfig>$`)
)

// CiscoAireOSDevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoAireOSDevice struct {
	IP                string
	Name              string
	Vendor            string
	Platform          string
	Connector         string
	UserPromptRE      *regexp.Regexp
	SuperUserPromptRE *regexp.Regexp
	ConfigPromtRE     *regexp.Regexp
	data.Credentials
	network.SSHParams
	network.SSHConn
	data.Variables
}

func (cd *CiscoAireOSDevice) ConnectWithSSH() error {

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

	cd.SendCommandWithSSH("config paging disable")

	return nil
}

func (cd CiscoAireOSDevice) DisconnectSSH() error {
	return cd.SSHConn.Session.Close()
}

func (cd CiscoAireOSDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandWithSSH(cd.SSHConn, command, cd.SuperUserPromptRE, 5)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (cd CiscoAireOSDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = cd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandsWithSSH(cd.SSHConn, commands, cd.SuperUserPromptRE, 5)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewCiscoAireOSDevice takes a NetDevice and initializes
// a CiscoAireOSDevice.
func NewCiscoAireOSDevice(nd NetDevice) CiscoAireOSDevice {
	cd := CiscoAireOSDevice{}
	cd.IP = nd.IP
	cd.Name = nd.Name
	cd.Vendor = nd.Vendor
	cd.Platform = nd.Platform
	cd.Connector = nd.Connector
	cd.Credentials = nd.Credentials
	cd.SSHParams = nd.SSHParams
	cd.Variables = nd.Variables

	// Prompts
	cd.UserPromptRE = CiscoAireOSUserPromptRE
	cd.SuperUserPromptRE = CiscoAireOSSuperUserPromptRE
	cd.ConfigPromtRE = CiscoAireOSConfigPromptRE

	// SSH Params
	network.InitSSHParams(&cd.SSHParams)

	return cd
}
