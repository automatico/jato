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
	IP             string `json:"ip"`
	Name           string `json:"name"`
	Vendor         string `json:"vendor"`
	Platform       string `json:"platform"`
	Connector      string `json:"connector"`
	data.Variables `json:"variables"`
	data.Credentials

	network.SSHParams
	network.TelnetParams
}
