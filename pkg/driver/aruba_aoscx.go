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
	ArubaCXUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}>\s$`)
	ArubaCXSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}#\s$`)
	ArubaCXConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9\.-]{1,31}\(config[a-z0-9-]{0,63}\)#\s$`)
)

// ArubaAOSCXDevice implements the TelnetDevice
// and SSHDevice interfaces
type ArubaAOSCXDevice struct {
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

func (ad *ArubaAOSCXDevice) ConnectWithSSH() error {

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

	ad.SendCommandWithSSH("no page")

	return nil
}

func (ad ArubaAOSCXDevice) DisconnectSSH() error {
	return ad.SSHConn.Session.Close()
}

func (ad ArubaAOSCXDevice) SendCommandWithSSH(command string) data.Result {

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

func (ad ArubaAOSCXDevice) SendCommandsWithSSH(commands []string) data.Result {

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

// NewArubaAOSCXDevice takes a NetDevice and initializes
// a ArubaAOSCXDevice.
func NewArubaAOSCXDevice(nd NetDevice) ArubaAOSCXDevice {
	ad := ArubaAOSCXDevice{}
	ad.IP = nd.IP
	ad.Name = nd.Name
	ad.Vendor = nd.Vendor
	ad.Platform = nd.Platform
	ad.Connector = nd.Connector
	ad.Credentials = nd.Credentials
	ad.SSHParams = nd.SSHParams

	// Prompts
	ad.UserPromptRE = ArubaCXUserPromptRE
	ad.SuperUserPromptRE = ArubaCXSuperUserPromptRE
	ad.ConfigPromtRE = ArubaCXConfigPromptRE

	// SSH Params
	if ad.SSHParams.Port == 0 {
		ad.SSHParams.Port = constant.SSHPort
	}
	if !ad.SSHParams.InsecureConnection {
		ad.SSHParams.InsecureConnection = true
	}
	if !ad.SSHParams.InsecureCyphers {
		ad.SSHParams.InsecureCyphers = false
	}
	if !ad.SSHParams.InsecureKeyExchange {
		ad.SSHParams.InsecureKeyExchange = false
	}
	return ad
}
