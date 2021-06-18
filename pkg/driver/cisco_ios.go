package driver

import (
	"fmt"
	"regexp"
	"time"

	"github.com/automatico/jato/internal/logger"
	"github.com/automatico/jato/pkg/constant"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
	"github.com/reiver/go-telnet"
)

var (
	CiscoUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\\-_@()/:]{1,63}>$`)
	CiscoSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\\-_@()/:]{1,63}#$`)
	CiscoConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}\([a-z0-9.\-@/:\+]{0,32}\)#$`)
)

// CiscoIOSDevice implements the TelnetDevice
// and SSHDevice interfaces
type CiscoIOSDevice struct {
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
	network.TelnetParams
	network.SSHConn
	TelnetConn *telnet.Conn
	data.Variables
}

func (d *CiscoIOSDevice) ConnectWithTelnet() error {

	conn, err := telnet.DialTo(fmt.Sprintf("%s:%d", d.IP, d.TelnetParams.Port))
	if err != nil {
		return err
	}

	_, err = network.SendCommandWithTelnet(conn, d.Username, constant.PasswordRE, 1)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	_, err = network.SendCommandWithTelnet(conn, d.Password, d.SuperUserPromptRE, 1)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	d.TelnetConn = conn

	d.SendCommandWithTelnet("terminal length 0")
	d.SendCommandWithTelnet("terminal width 0")

	return nil
}

func (d CiscoIOSDevice) DisconnectTelnet() error {
	return d.TelnetConn.Close()
}

func (d CiscoIOSDevice) SendCommandWithTelnet(cmd string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandWithTelnet(d.TelnetConn, cmd, d.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (d CiscoIOSDevice) SendCommandsWithTelnet(commands []string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandsWithTelnet(d.TelnetConn, commands, d.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

func (d *CiscoIOSDevice) ConnectWithSSH() error {

	clientConfig := network.SSHClientConfig(d.Credentials, d.SSHParams)

	sshConn, err := network.ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)
	if err != nil {
		return err
	}

	network.ReadSSH(sshConn.StdOut, d.SuperUserPromptRE, 2)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("terminal length 0")
	d.SendCommandWithSSH("terminal width 0")

	return nil
}

func (d CiscoIOSDevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d CiscoIOSDevice) SendCommandWithSSH(command string) data.Result {

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

func (d CiscoIOSDevice) SendCommandsWithSSH(commands []string) data.Result {

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

// NewCiscoIOSDevice takes a NetDevice and initializes
// a CiscoIOSDevice.
func NewCiscoIOSDevice(nd NetDevice) CiscoIOSDevice {
	d := CiscoIOSDevice{}
	d.IP = nd.IP
	d.Name = nd.Name
	d.Vendor = nd.Vendor
	d.Platform = nd.Platform
	d.Connector = nd.Connector
	d.Credentials = nd.Credentials
	d.SSHParams = nd.SSHParams
	d.TelnetParams = nd.TelnetParams
	d.Variables = nd.Variables

	// Prompts
	d.UserPromptRE = CiscoUserPromptRE
	d.SuperUserPromptRE = CiscoSuperUserPromptRE
	d.ConfigPromtRE = CiscoConfigPromptRE

	// SSH Params
	network.InitSSHParams(&d.SSHParams)

	// Telnet Params
	network.InitTelnetParams(&d.TelnetParams)

	return d
}
