package driver

import (
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
)

var (
	CiscoXRUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}#\s?$`)
	CiscoXRSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}#\s?$`)
	CiscoXRConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#$`)
)

// CiscoIOSXRDevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoIOSXRDevice struct {
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

func (cd *CiscoIOSXRDevice) ConnectWithSSH() error {

	clientConfig := network.SSHClientConfig(
		cd.Credentials.Username,
		cd.Credentials.Password,
		cd.SSHParams.InsecureConnection,
		cd.SSHParams.InsecureCyphers,
		cd.SSHParams.InsecureKeyExchange,
	)

	sshConn := network.ConnectWithSSH(cd.IP, cd.SSHParams.Port, clientConfig)

	network.ReadSSH(sshConn.StdOut, cd.SuperUserPromptRE, 2)

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
		result.Error = err
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
		result.Error = err
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
	cd.Variables = nd.Variables

	// Prompts
	cd.UserPromptRE = CiscoXRUserPromptRE
	cd.SuperUserPromptRE = CiscoXRSuperUserPromptRE
	cd.ConfigPromtRE = CiscoXRConfigPromptRE

	// SSH Params
	network.InitSSHParams(&cd.SSHParams)

	return cd
}
