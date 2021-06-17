package driver

import (
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
)

var (
	ArubaCXUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}>\s$`)
	ArubaCXSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}#\s$`)
	ArubaCXConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}\(config[a-z0-9-]{0,63}\)#\s$`)
)

// ArubaAOSCXDevice implements the TelnetDevice
// and SSHDevice interfaces
type ArubaAOSCXDevice struct {
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

func (d *ArubaAOSCXDevice) ConnectWithSSH() error {

	clientConfig := network.SSHClientConfig(d.Credentials, d.SSHParams)

	sshConn := network.ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)

	network.ReadSSH(sshConn.StdOut, d.SuperUserPromptRE, 2)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("no page")

	return nil
}

func (d ArubaAOSCXDevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d ArubaAOSCXDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandWithSSH(d.SSHConn, command, d.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (d ArubaAOSCXDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandsWithSSH(d.SSHConn, commands, d.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewArubaAOSCXDevice takes a NetDevice and initializes
// a ArubaAOSCXDevice.
func NewArubaAOSCXDevice(nd NetDevice) ArubaAOSCXDevice {
	d := ArubaAOSCXDevice{}
	d.IP = nd.IP
	d.Name = nd.Name
	d.Vendor = nd.Vendor
	d.Platform = nd.Platform
	d.Connector = nd.Connector
	d.Credentials = nd.Credentials
	d.SSHParams = nd.SSHParams
	d.Variables = nd.Variables

	// Prompts
	d.UserPromptRE = ArubaCXUserPromptRE
	d.SuperUserPromptRE = ArubaCXSuperUserPromptRE
	d.ConfigPromtRE = ArubaCXConfigPromptRE

	// SSH Params
	network.InitSSHParams(&d.SSHParams)

	return d
}
