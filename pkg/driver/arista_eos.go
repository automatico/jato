package driver

import (
	"regexp"
)

// NewAristaEOSDevice takes a NetDevice and initializes
// a AristaEOSDevice.
func NewAristaEOSDevice(d NetDevice) NetDevice {

	// Prompts
	d.UserPromptRE = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,63}>$`)
	d.SuperUserPromptRE = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,63}#$`)
	d.ConfigPromtRE = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,63}\(config[a-z0-9-]{0,63}\)#$`)

	// SSH Params
	InitSSHParams(&d.SSHParams)

	// Timeout
	d.Timeout = 5

	return d
}

func AristaEOSConnectWithSSH(d *NetDevice) error {

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
	d.SendCommandWithSSH("terminal width 32767")

	return nil
}
