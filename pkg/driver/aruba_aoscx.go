package driver

import (
	"regexp"
)

// NewArubaAOSCXDevice takes a NetDevice and initializes
// a ArubaAOSCXDevice.
func NewArubaAOSCXDevice(d NetDevice) NetDevice {

	// Prompts
	d.UserPromptRE = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}>\s$`)
	d.SuperUserPromptRE = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}#\s$`)
	d.ConfigPromtRE = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}\(config[a-z0-9-]{0,63}\)#\s$`)

	// SSH Params
	InitSSHParams(&d.SSHParams)

	// Timeout
	d.Timeout = 5

	return d
}

func ArubaAOSCXConnectWithSSH(d *NetDevice) error {

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

	d.SendCommandWithSSH("no page")

	return nil
}
