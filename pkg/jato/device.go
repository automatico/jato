package jato

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

// Device represents a managed device
type Device struct {
	IP        string `json:"ip"`
	Name      string `json:"name"`
	Vendor    string `json:"vendor"`
	Platform  string `json:"platform"`
	Connector string `json:"connector"`
}

// Devices holds a collection of Device structs
type Devices struct {
	Devices []Device `json:"devices"`
}

type TelnetParams struct {
	Port int
}

type SSHParams struct {
	Port               int
	InsecureConnection bool
	InsecureCyphers    bool
}

type NetDevice struct {
	IP        string `json:"ip"`
	Name      string `json:"name"`
	Vendor    string `json:"vendor"`
	Platform  string `json:"platform"`
	Connector string `json:"connector"`
	Prompt    string
	Credentials
	SSHParams
	TelnetParams
}

func (nd NetDevice) ConnectWithTelnet() net.Conn {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", nd.IP, nd.TelnetParams.Port))
	if err != nil {
		fmt.Println("dial error:", err)
	}
	commands := CommandExpect{
		[]Expect{
			{Command: "", Expecting: "Username:", Timeout: 2},
			{Command: nd.Credentials.Username, Expecting: "Password:", Timeout: 2},
			{Command: nd.Credentials.Password, Expecting: nd.Prompt, Timeout: 2},
		},
	}
	for _, cmd := range commands.CommandExpect {
		result, err := TelnetExpecter(conn, cmd.Command, cmd.Expecting, cmd.Timeout)
		if err != nil {
			fmt.Println(result)
			fmt.Println(err)
		}
	}
	return conn
}

func (nd NetDevice) ConnectWithSSH() SSHConnection {

	sshConn := SSHConnection{}

	sshClientConfig := SSHClientConfig(
		nd.Credentials.Username,
		nd.Credentials.Password,
		nd.SSHParams.InsecureConnection,
		nd.SSHParams.InsecureCyphers,
	)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", nd.IP, nd.SSHParams.Port), sshClientConfig)
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
	if err := session.RequestPty("xterm", 0, 200, modes); err != nil {
		session.Close()
		fmt.Println(err)
	}

	if err := session.Shell(); err != nil {
		session.Close()
		fmt.Println(err)
	}
	readBuff(nd.Prompt, stdOut, 2)

	sshConn.Session = session
	sshConn.StdIn = stdIn
	sshConn.StdOut = stdOut

	return sshConn
}
