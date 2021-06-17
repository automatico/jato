package driver

import (
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
)

var (
	AristaUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,63}>$`)
	AristaSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,63}#$`)
	AristaConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,63}\(config[a-z0-9-]{0,63}\)#$`)
)

// AristaEOSDevice implements the TelnetDevice
// and SSHDevice interfaces
type AristaEOSDevice struct {
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

func (d *AristaEOSDevice) ConnectWithSSH() error {

	clientConfig := network.SSHClientConfig(d.Credentials, d.SSHParams)

	sshConn := network.ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)

	network.ReadSSH(sshConn.StdOut, d.SuperUserPromptRE, 2)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("terminal length 0")
	d.SendCommandWithSSH("terminal width 32767")

	return nil
}

func (d AristaEOSDevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d AristaEOSDevice) SendCommandWithSSH(command string) data.Result {

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

func (d AristaEOSDevice) SendCommandsWithSSH(commands []string) data.Result {

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

// NewAristaEOSDevice takes a NetDevice and initializes
// a AristaEOSDevice.
func NewAristaEOSDevice(nd NetDevice) AristaEOSDevice {
	d := AristaEOSDevice{}
	d.IP = nd.IP
	d.Name = nd.Name
	d.Vendor = nd.Vendor
	d.Platform = nd.Platform
	d.Connector = nd.Connector
	d.Credentials = nd.Credentials
	d.SSHParams = nd.SSHParams
	d.Variables = nd.Variables

	// Prompts
	d.UserPromptRE = AristaUserPromptRE
	d.SuperUserPromptRE = AristaSuperUserPromptRE
	d.ConfigPromtRE = AristaConfigPromptRE

	// SSH Params
	network.InitSSHParams(&d.SSHParams)

	return d
}
