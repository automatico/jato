package driver

import (
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
	return cd
}
