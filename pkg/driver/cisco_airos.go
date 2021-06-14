package driver

import (
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
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

func (d *CiscoAireOSDevice) ConnectWithSSH() error {

	clientConfig := network.SSHClientConfig(
		d.Credentials.Username,
		d.Credentials.Password,
		d.SSHParams.InsecureConnection,
		d.SSHParams.InsecureCyphers,
		d.SSHParams.InsecureKeyExchange,
	)

	sshConn := network.ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)

	network.ReadSSH(sshConn.StdOut, d.SuperUserPromptRE, 2)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("config paging disable")

	return nil
}

func (d CiscoAireOSDevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d CiscoAireOSDevice) SendCommandWithSSH(command string) data.Result {

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

func (d CiscoAireOSDevice) SendCommandsWithSSH(commands []string) data.Result {

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

// NewCiscoAireOSDevice takes a NetDevice and initializes
// a CiscoAireOSDevice.
func NewCiscoAireOSDevice(nd NetDevice) CiscoAireOSDevice {
	d := CiscoAireOSDevice{}
	d.IP = nd.IP
	d.Name = nd.Name
	d.Vendor = nd.Vendor
	d.Platform = nd.Platform
	d.Connector = nd.Connector
	d.Credentials = nd.Credentials
	d.SSHParams = nd.SSHParams
	d.Variables = nd.Variables

	// Prompts
	d.UserPromptRE = CiscoAireOSUserPromptRE
	d.SuperUserPromptRE = CiscoAireOSSuperUserPromptRE
	d.ConfigPromtRE = CiscoAireOSConfigPromptRE

	// SSH Params
	network.InitSSHParams(&d.SSHParams)

	return d
}
