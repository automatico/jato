package driver

import (
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
)

var (
	CiscoASAUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\\-_@()/:]{1,63}>\s$`)
	CiscoASASuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}#\s$`)
	CiscoASAConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#\s$`)
)

// CiscoASADevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoASADevice struct {
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

func (d *CiscoASADevice) ConnectWithSSH() error {

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

	d.SendCommandWithSSH("terminal pager 0")

	return nil
}

func (d CiscoASADevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d CiscoASADevice) SendCommandWithSSH(command string) data.Result {

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

func (d CiscoASADevice) SendCommandsWithSSH(commands []string) data.Result {

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

// NewCiscoASADevice takes a NetDevice and initializes
// a CiscoASADevice.
func NewCiscoASADevice(nd NetDevice) CiscoASADevice {
	d := CiscoASADevice{}
	d.IP = nd.IP
	d.Name = nd.Name
	d.Vendor = nd.Vendor
	d.Platform = nd.Platform
	d.Connector = nd.Connector
	d.Credentials = nd.Credentials
	d.SSHParams = nd.SSHParams
	d.Variables = nd.Variables

	// Prompts
	d.UserPromptRE = CiscoASAUserPromptRE
	d.SuperUserPromptRE = CiscoASASuperUserPromptRE
	d.ConfigPromtRE = CiscoASAConfigPromptRE

	// SSH Params
	network.InitSSHParams(&d.SSHParams)

	return d
}
