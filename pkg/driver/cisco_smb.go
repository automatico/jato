package driver

import (
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
)

var (
	CiscoSMBUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\\-_@()/:]{1,63}>$`)
	CiscoSMBSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\\-_@()/:]{1,63}#$`)
	CiscoSMBConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}\([a-z0-9.\-@/:\+]{0,32}\)#$`)
)

// CiscoSMBDevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoSMBDevice struct {
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

func (d *CiscoSMBDevice) ConnectWithSSH() error {

	clientConfig := network.SSHClientConfig(d.Credentials, d.SSHParams)

	sshConn := network.ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)

	network.ReadSSH(sshConn.StdOut, d.SuperUserPromptRE, 5)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("terminal datadump")
	d.SendCommandWithSSH("terminal width 512")

	return nil
}

func (d CiscoSMBDevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d CiscoSMBDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	// Cisco SMB devices are really slow to output to the terminal.
	cmdOut, err := network.SendCommandWithSSH(d.SSHConn, command, d.SuperUserPromptRE, 120)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (d CiscoSMBDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	// Cisco SMB devices are really slow to output to the terminal.
	cmdOut, err := network.SendCommandsWithSSH(d.SSHConn, commands, d.SuperUserPromptRE, 120)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewCiscoSMBDevice takes a NetDevice and initializes
// a CiscoSMBDevice.
func NewCiscoSMBDevice(nd NetDevice) CiscoSMBDevice {
	d := CiscoSMBDevice{}
	d.IP = nd.IP
	d.Name = nd.Name
	d.Vendor = nd.Vendor
	d.Platform = nd.Platform
	d.Connector = nd.Connector
	d.Credentials = nd.Credentials
	d.SSHParams = nd.SSHParams
	d.Variables = nd.Variables

	// Prompts
	d.UserPromptRE = CiscoSMBUserPromptRE
	d.SuperUserPromptRE = CiscoSMBSuperUserPromptRE
	d.ConfigPromtRE = CiscoSMBConfigPromptRE

	// SSH Params
	network.InitSSHParams(&d.SSHParams)

	return d
}
