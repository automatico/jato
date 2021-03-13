package jato

type Connector interface {
	Auth()
	Connect()
}

type Jato struct {
	UserCredentials
	Devices
	CommandExpect
}
