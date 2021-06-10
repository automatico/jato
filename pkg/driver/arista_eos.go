package driver

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/constant"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
	"golang.org/x/crypto/ssh"
)

var (
	AristaUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.-]{1,63}>$`)
	AristaSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.-]{1,63}#$`)
	AristaConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.-]{1,63}\(config[a-z0-9-]{0,63}\)#$`)
)

// AristaEOSDevice implements the TelnetDevice
// and SSHDevice interfaces
type AristaEOSDevice struct {
	IP                string `json:"ip"`
	Name              string `json:"name"`
	Vendor            string `json:"vendor"`
	Platform          string `json:"platform"`
	Connector         string `json:"connector"`
	UserPromptRE      *regexp.Regexp
	SuperUserPromptRE *regexp.Regexp
	ConfigPromtRE     *regexp.Regexp
	data.Credentials
	network.SSHParams
	network.SSHConn
}

func (ad *AristaEOSDevice) ConnectWithSSH() error {

	sshConn := network.SSHConn{}

	clientConfig := network.SSHClientConfig(
		ad.Credentials.Username,
		ad.Credentials.Password,
		ad.SSHParams.InsecureConnection,
		ad.SSHParams.InsecureCyphers,
		ad.SSHParams.InsecureKeyExchange,
	)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ad.IP, ad.SSHParams.Port), clientConfig)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		fmt.Println(err)
	}

	stdOut, err := session.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	stdIn, err := session.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = session.RequestPty("xterm", 0, 200, modes)
	if err != nil {
		session.Close()
		fmt.Println(err)
	}

	err = session.Shell()
	if err != nil {
		session.Close()
		fmt.Println(err)
	}

	network.ReadSSH(stdOut, ad.SuperUserPromptRE, 2)

	sshConn.Session = session
	sshConn.StdIn = stdIn
	sshConn.StdOut = stdOut

	ad.SSHConn = sshConn

	ad.SendCommandWithSSH("terminal length 0")
	ad.SendCommandWithSSH("terminal width 32767")

	return nil
}

func (ad AristaEOSDevice) DisconnectSSH() error {
	return ad.SSHConn.Session.Close()
}

func (ad AristaEOSDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = ad.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandWithSSH(ad.SSHConn, command, ad.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (ad AristaEOSDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = ad.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandsWithSSH(ad.SSHConn, commands, ad.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewAristaEOSDevice takes a NetDevice and initializes
// a AristaEOSDevice.
func NewAristaEOSDevice(nd NetDevice) AristaEOSDevice {
	ad := AristaEOSDevice{}
	ad.IP = nd.IP
	ad.Name = nd.Name
	ad.Vendor = nd.Vendor
	ad.Platform = nd.Platform
	ad.Connector = nd.Connector
	ad.Credentials = nd.Credentials
	ad.SSHParams = nd.SSHParams

	// Prompts
	ad.UserPromptRE = AristaUserPromptRE
	ad.SuperUserPromptRE = AristaSuperUserPromptRE
	ad.ConfigPromtRE = AristaConfigPromptRE

	// SSH Params
	if ad.SSHParams.Port == 0 {
		ad.SSHParams.Port = constant.SSHPort
	}
	if !ad.SSHParams.InsecureConnection {
		ad.SSHParams.InsecureConnection = true
	}
	if !ad.SSHParams.InsecureCyphers {
		ad.SSHParams.InsecureCyphers = true
	}
	if !ad.SSHParams.InsecureKeyExchange {
		ad.SSHParams.InsecureKeyExchange = true
	}
	return ad
}
