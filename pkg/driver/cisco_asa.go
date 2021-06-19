package driver

import (
	"regexp"
)

// NewCiscoASADevice takes a NetDevice and initializes
// a CiscoASADevice.
func NewCiscoASADevice(d NetDevice) NetDevice {

	// Prompts
	d.UserPromptRE = regexp.MustCompile(`(?im)[a-z0-9\-]{1,63}>\s$`)
	d.SuperUserPromptRE = regexp.MustCompile(`(?im)[a-z0-9\-]{1,63}#\s$`)
	d.ConfigPromtRE = regexp.MustCompile(`(?im)[a-z0-9\-]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#\s$`)

	// SSH Params
	InitSSHParams(&d.SSHParams)

	// Timeout
	d.Timeout = 5

	return d
}

func CiscoASAConnectWithSSH(d *NetDevice) error {

	clientConfig, err := SSHClientConfig(d.Credentials, d.SSHParams)
	if err != nil {
		return err
	}

	sshConn, err := ConnectWithSSH(d.IP, d.SSHParams.Port, clientConfig)
	if err != nil {
		return err
	}

	ReadSSH(sshConn.StdOut, d.SuperUserPromptRE, 2)

	d.SSHConn = sshConn

	d.SendCommandWithSSH("terminal pager 0")

	return err
}
