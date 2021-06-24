package driver

import (
	"regexp"
)

// NewCiscoNXOSDevice takes a NetDevice and initializes
// a CiscoNXOSDevice.
func NewCiscoNXOSDevice(d NetDevice) NetDevice {
	// Prompts
	d.Prompt.User = regexp.MustCompile(`(?im)[a-z0-9.\\-_@()/:]{1,63}>\s$`)
	d.Prompt.SuperUser = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}#\s$`)
	d.Prompt.Config = regexp.MustCompile(`(?im)[a-z0-9.\-_@/:]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#\s$`)

	// SSH Params
	InitSSHParams(&d.SSHParams)

	// Timeout
	d.Timeout = 5

	return d
}

func CiscoNXOSConnectWithSSH(d *NetDevice) error {

	clientConfig, err := SSHClientConfig(d.Credentials, d.SSHParams)
	if err != nil {
		return err
	}

	sshConn, err := ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)
	if err != nil {
		return err
	}

	ReadSSH(sshConn.StdOut, d.Prompt.SuperUser, 2)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("terminal length 0")
	d.SendCommandWithSSH("terminal width 511")

	return nil
}
