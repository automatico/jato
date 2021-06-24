package driver

import (
	"fmt"
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/reiver/go-telnet"
)

type Prompt struct {
	User      *regexp.Regexp
	SuperUser *regexp.Regexp
	Config    *regexp.Regexp
}

// Devices holds a collection of NetDevice structs
type NetDevices struct {
	Devices []NetDevice `json:"devices"`
}

type NetDevice struct {
	IP             string `json:"ip"`
	Name           string `json:"name"`
	Vendor         string `json:"vendor"`
	Platform       string `json:"platform"`
	Connector      string `json:"connector"`
	SSHParams      `json:"sshParams"`
	TelnetParams   `json:"telnetParams"`
	data.Variables `json:"variables"`
	Timeout        int64
	Prompt         Prompt
	SSHConn
	TelnetConn *telnet.Conn
	data.Credentials
}

func (d *NetDevice) ConnectWithSSH() error {

	vendorPlatform := fmt.Sprintf("%s_%s", d.Vendor, d.Platform)
	switch vendorPlatform {
	case "arista_eos":
		err := AristaEOSConnectWithSSH(d)
		if err != nil {
			return err
		}
	case "aruba_aoscx":
		err := ArubaAOSCXConnectWithSSH(d)
		if err != nil {
			return err
		}
	case "cisco_aireos":
		err := CiscoAireOSConnectWithSSH(d)
		if err != nil {
			return err
		}
	case "cisco_asa":
		err := CiscoASAConnectWithSSH(d)
		if err != nil {
			return err
		}
	case "cisco_ios":
		err := CiscoIOSConnectWithSSH(d)
		if err != nil {
			return err
		}
	case "cisco_iosxr":
		err := CiscoIOSXRConnectWithSSH(d)
		if err != nil {
			return err
		}
	case "cisco_nxos":
		err := CiscoNXOSConnectWithSSH(d)
		if err != nil {
			return err
		}
	case "cisco_smb":
		err := CiscoSMBConnectWithSSH(d)
		if err != nil {
			return err
		}
	case "juniper_junos":
		err := JuniperJunosConnectWithSSH(d)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("device: %s with vendor: %s and platform: %s not supported", d.Name, d.Vendor, d.Platform)
	}
	return nil

}

func (d NetDevice) DisconnectSSH() error {
	return d.SSHConn.Session.Close()
}

func (d NetDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := SendCommandWithSSH(d.SSHConn, command, d.Prompt.SuperUser, d.Timeout)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (d NetDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := SendCommandsWithSSH(d.SSHConn, commands, d.Prompt.SuperUser, d.Timeout)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

func (d *NetDevice) ConnectWithTelnet() error {

	vendorPlatform := fmt.Sprintf("%s_%s", d.Vendor, d.Platform)
	switch vendorPlatform {
	case "cisco_ios":
		err := CiscoIOSConnectWithTelnet(d)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("device: %s with vendor: %s and platform: %s not supported", d.Name, d.Vendor, d.Platform)
	}
	return nil

}

func (d NetDevice) DisconnectTelnet() error {
	return d.TelnetConn.Close()
}

func (d NetDevice) SendCommandWithTelnet(cmd string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := SendCommandWithTelnet(d.TelnetConn, cmd, d.Prompt.SuperUser, d.Timeout)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (d NetDevice) SendCommandsWithTelnet(commands []string) data.Result {

	result := data.Result{}

	result.Device = d.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := SendCommandsWithTelnet(d.TelnetConn, commands, d.Prompt.SuperUser, d.Timeout)
	if err != nil {
		result.OK = false
		result.Error = err
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}
