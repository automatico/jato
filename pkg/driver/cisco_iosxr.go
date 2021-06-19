package driver

import (
	"regexp"
)

// NewCiscoIOSXRDevice takes a NetDevice and initializes
// a CiscoIOSXRDevice.
func NewCiscoIOSXRDevice(d NetDevice) NetDevice {

	// Prompts
	d.UserPromptRE = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}#\s?$`)
	d.SuperUserPromptRE = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}#\s?$`)
	d.ConfigPromtRE = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#$`)

	// SSH Params
	InitSSHParams(&d.SSHParams)

	// Timeout
	d.Timeout = 5

	return d
}

func CiscoIOSXRConnectWithSSH(d *NetDevice) error {

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

	d.SendCommandWithSSH("terminal length 0")
	d.SendCommandWithSSH("terminal width 0")

	return nil
}
