package jato

type Connector interface {
	Auth()
	Connect()
}

type NetDevice interface {
	Connect()
	DisablePaging()
}

type Jato struct {
	UserCredentials
	Devices
	CommandExpect
}
