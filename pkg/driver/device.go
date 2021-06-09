package driver

import (
	"github.com/automatico/jato/pkg/constant"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
)

// Devices holds a collection of Device structs
type Devices struct {
	Devices []NetDevice `json:"devices"`
}

type NetDevice struct {
	IP        string `json:"ip"`
	Name      string `json:"name"`
	Vendor    string `json:"vendor"`
	Platform  string `json:"platform"`
	Connector string `json:"connector"`
	data.Credentials
	network.SSHParams
	network.TelnetParams
}

func NewCiscoIOSDevice(nd NetDevice) CiscoIOSDevice {
	cd := CiscoIOSDevice{}
	cd.IP = nd.IP
	cd.Name = nd.Name
	cd.Vendor = nd.Vendor
	cd.Platform = nd.Platform
	cd.Connector = nd.Connector
	cd.Credentials = nd.Credentials
	cd.SSHParams = nd.SSHParams
	cd.TelnetParams = nd.TelnetParams

	// Prompts
	cd.UserPromptRE = CiscoUserPromptRE
	cd.SuperUserPromptRE = CiscoSuperUserPromptRE
	cd.ConfigPromtRE = CiscoConfigPromptRE

	// Paging
	cd.DisablePaging = CiscoDisablePaging

	// SSH Params
	if cd.SSHParams.Port == 0 {
		cd.SSHParams.Port = constant.SSHPort
	}
	if !cd.SSHParams.InsecureConnection {
		cd.SSHParams.InsecureConnection = true
	}
	if !cd.SSHParams.InsecureCyphers {
		cd.SSHParams.InsecureCyphers = true
	}

	// Telnet Params
	if cd.TelnetParams.Port == 0 {
		cd.TelnetParams.Port = constant.TelnetPort
	}
	return cd
}
