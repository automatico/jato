package driver

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/network"
	"golang.org/x/crypto/ssh"
)

var (
	JuniperUserPromptRE      *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@()/:]{1,63}>\s$`)
	JuniperSuperUserPromptRE *regexp.Regexp = regexp.MustCompile(`(?im)[a-z0-9.\-_@()/:]{1,63}>\s$`)
	JuniperConfigPromptRE    *regexp.Regexp = regexp.MustCompile(`(?im)(\[edit\]\n){0,1}[a-z0-9.\-_@()/:]{1,63}#\s?$`)
)

// JuniperJunosDevice implements the TelnetDevice
// and SSHDevice interfaces
type JuniperJunosDevice struct {
	IP                string
	Name              string
	Vendor            string
	Platform          string
	Connector         string
	UserPromptRE      *regexp.Regexp
	SuperUserPromptRE *regexp.Regexp
	ConfigPromtRE     *regexp.Regexp
	data.Credentials
	network.SSHParams
	network.SSHConn
	data.Variables
}

func (jd *JuniperJunosDevice) ConnectWithSSH() error {

	sshConn := network.SSHConn{}

	clientConfig := network.SSHClientConfig(
		jd.Credentials.Username,
		jd.Credentials.Password,
		jd.SSHParams.InsecureConnection,
		jd.SSHParams.InsecureCyphers,
		jd.SSHParams.InsecureKeyExchange,
	)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", jd.IP, jd.SSHParams.Port), clientConfig)
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

	network.ReadSSH(stdOut, jd.SuperUserPromptRE, 2)

	sshConn.Session = session
	sshConn.StdIn = stdIn
	sshConn.StdOut = stdOut

	jd.SSHConn = sshConn

	jd.SendCommandWithSSH("set cli screen-length 0")
	jd.SendCommandWithSSH("set cli screen-width 0")

	return nil
}

func (jd JuniperJunosDevice) DisconnectSSH() error {
	return jd.SSHConn.Session.Close()
}

func (jd JuniperJunosDevice) SendCommandWithSSH(command string) data.Result {

	result := data.Result{}

	result.Device = jd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandWithSSH(jd.SSHConn, command, jd.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = append(result.CommandOutputs, cmdOut)
	result.OK = true
	return result
}

func (jd JuniperJunosDevice) SendCommandsWithSSH(commands []string) data.Result {

	result := data.Result{}

	result.Device = jd.Name
	result.Timestamp = time.Now().Unix()

	cmdOut, err := network.SendCommandsWithSSH(jd.SSHConn, commands, jd.SuperUserPromptRE, 2)
	if err != nil {
		result.OK = false
		return result
	}

	result.CommandOutputs = cmdOut
	result.OK = true
	return result
}

// NewJuniperJunosDevice takes a NetDevice and initializes
// a JuniperJunosDevice.
func NewJuniperJunosDevice(nd NetDevice) JuniperJunosDevice {
	jd := JuniperJunosDevice{}
	jd.IP = nd.IP
	jd.Name = nd.Name
	jd.Vendor = nd.Vendor
	jd.Platform = nd.Platform
	jd.Connector = nd.Connector
	jd.Credentials = nd.Credentials
	jd.SSHParams = nd.SSHParams
	jd.Variables = nd.Variables

	// Prompts
	jd.UserPromptRE = JuniperUserPromptRE
	jd.SuperUserPromptRE = JuniperSuperUserPromptRE
	jd.ConfigPromtRE = JuniperConfigPromptRE

	// SSH Params
	network.InitSSHParams(&jd.SSHParams)

	return jd
}
