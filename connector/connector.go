package connector

import (
	"github.com/automatico/jato/credentials"
	"github.com/automatico/jato/device"
	"github.com/automatico/jato/expecter"
)

type Connector interface {
	Auth()
	Connect()
}

type Jato struct {
	credentials.UserCredentials
	device.Devices
	expecter.CommandExpect
}
