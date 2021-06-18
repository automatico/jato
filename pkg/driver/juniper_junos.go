package driver

import (
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
)

var (
	JuniperUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@()/:]{1,63}>\s$`)
	JuniperSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@()/:]{1,63}>\s$`)
	JuniperConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)(\[edit\]\n){0,1}[a-z0-9.\-_@()/:]{1,63}#\s?$`)
)

// JuniperJunosDevice implements the TelnetDevice
// and SSHDevice interfaces
type JuniperJunosDevice struct {
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

func (d *JuniperJunosDevice) ConnectWithSSH() error {

	clientConfig := network.SSHClientConfig(d.Credentials, d.SSHParams)

	sshConn, err := network.ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)
	if err != nil {
		return err
	}

	network.ReadSSH(sshConn.StdOut, d.SuperUserPromptRE, 2)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("set cli screen-length 0")
	d.SendCommandWithSSH("set cli screen-width 0")

	return nil
}

func (d JuniperJunosDevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d JuniperJunosDevice) SendCommandWithSSH(command string) data.Result {

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

func (d JuniperJunosDevice) SendCommandsWithSSH(commands []string) data.Result {

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

// NewJuniperJunosDevice takes a NetDevice and initializes
// a JuniperJunosDevice.
func NewJuniperJunosDevice(nd NetDevice) JuniperJunosDevice {
	d := JuniperJunosDevice{}
	d.IP = nd.IP
	d.Name = nd.Name
	d.Vendor = nd.Vendor
	d.Platform = nd.Platform
	d.Connector = nd.Connector
	d.Credentials = nd.Credentials
	d.SSHParams = nd.SSHParams
	d.Variables = nd.Variables

	// Prompts
	d.UserPromptRE = JuniperUserPromptRE
	d.SuperUserPromptRE = JuniperSuperUserPromptRE
	d.ConfigPromtRE = JuniperConfigPromptRE

	// SSH Params
	network.InitSSHParams(&d.SSHParams)

	return d
}
