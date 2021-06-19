package driver

import (
	"regexp"
)

// NewCiscoAireOSDevice takes a NetDevice and initializes
// a CiscoAireOSDevice.
func NewCiscoAireOSDevice(d NetDevice) NetDevice {

	// Prompts
	d.UserPromptRE = regexp.MustCompile(`(?im)^\([a-z0-9.\\-_\s@()/:]{1,63}\)\s>$`)
	d.SuperUserPromptRE = regexp.MustCompile(`(?im)^\([a-z0-9.\\-_\s@()/:]{1,63}\)\s>$`)
	d.ConfigPromtRE = regexp.MustCompile(`(?im)^\([a-z0-9.\\-_\s@()/:]{1,63}\)\sconfig>$`)

	// SSH Params
	InitSSHParams(&d.SSHParams)

	// Timeout
	d.Timeout = 5

	return d
}

func CiscoAireOSConnectWithSSH(d *NetDevice) error {

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

	d.SendCommandWithSSH("config paging disable")

	return nil
}
