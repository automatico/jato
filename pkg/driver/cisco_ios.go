package driver

import (
	"fmt"
	"regexp"

	"github.com/automatico/jato/internal/logger"
	"github.com/automatico/jato/pkg/constant"
	"github.com/reiver/go-telnet"
)

// NewCiscoIOSDevice takes a NetDevice and initializes
// a CiscoIOSDevice.
func NewCiscoIOSDevice(d NetDevice) NetDevice {
	// Prompts
	d.Prompt.User = regexp.MustCompile(`(?im)^[a-z0-9.\\-_@()/:]{1,63}>$`)
	d.Prompt.SuperUser = regexp.MustCompile(`(?im)^[a-z0-9.\\-_@()/:]{1,63}#$`)
	d.Prompt.Config = regexp.MustCompile(`(?im)^[a-z0-9.\-_@/:]{1,63}\([a-z0-9.\-@/:\+]{0,32}\)#$`)

	// SSH Params
	InitSSHParams(&d.SSHParams)

	// Telnet Params
	InitTelnetParams(&d.TelnetParams)

	// Timeout
	d.Timeout = 5

	return d
}

func CiscoIOSConnectWithSSH(d *NetDevice) error {

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
	d.SendCommandWithSSH("terminal width 0")

	return nil
}

func CiscoIOSConnectWithTelnet(d *NetDevice) error {

	conn, err := telnet.DialTo(fmt.Sprintf("%s:%d", d.IP, d.TelnetParams.Port))
	if err != nil {
		return err
	}

	_, err = SendCommandWithTelnet(conn, d.Username, constant.PasswordRE, 2)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	_, err = SendCommandWithTelnet(conn, d.Password, d.Prompt.SuperUser, 2)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	d.TelnetConn = conn

	d.SendCommandWithTelnet("terminal length 0")
	d.SendCommandWithTelnet("terminal width 0")

	return nil
}
