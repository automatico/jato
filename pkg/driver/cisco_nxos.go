package driver

import (
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
)

var (
	CiscoNXOSUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\\-_@()/:]{1,63}>\s$`)
	CiscoNXOSSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}#\s$`)
	CiscoNXOSConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#\s$`)
)

// CiscoNXOSDevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoNXOSDevice struct {
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

func (d CiscoNXOSDevice) GetName() string {
	return d.Name
}

func (d *CiscoNXOSDevice) ConnectWithSSH() error {

	clientConfig, err := network.SSHClientConfig(d.Credentials, d.SSHParams)
	if err != nil {
		return err
	}

	sshConn, err := network.ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)
	if err != nil {
		return err
	}

	network.ReadSSH(sshConn.StdOut, d.SuperUserPromptRE, 2)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("terminal length 0")
	d.SendCommandWithSSH("terminal width 511")

	return nil
}

func (d CiscoNXOSDevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d CiscoNXOSDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandWithSSH(d.SSHConn, command, d.SuperUserPromptRE, 5)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (d CiscoNXOSDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandsWithSSH(d.SSHConn, commands, d.SuperUserPromptRE, 5)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewCiscoNXOSDevice takes a NetDevice and initializes
// a CiscoNXOSDevice.
func NewCiscoNXOSDevice(nd NetDevice) CiscoNXOSDevice {
	d := CiscoNXOSDevice{}
	d.IP = nd.IP
	d.Name = nd.Name
	d.Vendor = nd.Vendor
	d.Platform = nd.Platform
	d.Connector = nd.Connector
	d.Credentials = nd.Credentials
	d.SSHParams = nd.SSHParams
	d.Variables = nd.Variables

	// Prompts
	d.UserPromptRE = CiscoNXOSUserPromptRE
	d.SuperUserPromptRE = CiscoNXOSSuperUserPromptRE
	d.ConfigPromtRE = CiscoNXOSConfigPromptRE

	// SSH Params
	network.InitSSHParams(&d.SSHParams)

	return d
}
