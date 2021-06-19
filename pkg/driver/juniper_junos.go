package driver

import (
	"regexp"
)

// NewJuniperJunosDevice takes a NetDevice and initializes
// a JuniperJunosDevice.
func NewJuniperJunosDevice(d NetDevice) NetDevice {
	// Prompts
	d.UserPromptRE = regexp.MustCompile(`(?im)[a-z0-9.\-_@()/:]{1,63}>\s$`)
	d.SuperUserPromptRE = regexp.MustCompile(`(?im)[a-z0-9.\-_@()/:]{1,63}>\s$`)
	d.ConfigPromtRE = regexp.MustCompile(`(?im)(\[edit\]\n){0,1}[a-z0-9.\-_@()/:]{1,63}#\s?$`)

	// SSH Params
	InitSSHParams(&d.SSHParams)

	// Timeout
	d.Timeout = 5

	return d
}

func JuniperJunosConnectWithSSH(d *NetDevice) error {

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

	d.SendCommandWithSSH("set cli screen-length 0")
	d.SendCommandWithSSH("set cli screen-width 0")

	return nil
}
