package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Commands to run against a device
type Commands struct {
	Commands []string `json:"commands"`
}

// Expect like interface
func expecter(cmd string, expect string, timeout int, sshIn io.WriteCloser, sshOut io.Reader) string {
	if _, err := writeBuff(cmd, sshIn); err != nil {
		handleError(err, true, "Failed to run: %s")
	}
	result := readBuff(expect, sshOut, timeout)

	return result
}

func readBuffForString(whattoexpect string, sshOut io.Reader, buffRead chan<- string) {
	buf := make([]byte, 1000)
	n, err := sshOut.Read(buf) //this reads the ssh terminal
	waitingString := ""
	if err == nil {
		waitingString = string(buf[:n])
	}
	for (err == nil) && (!strings.Contains(waitingString, whattoexpect)) {
		n, err = sshOut.Read(buf)
		waitingString += string(buf[:n])
		// fmt.Println(waitingString) //uncommenting this might help you debug if you are coming into errors with timeouts when correct details entered

	}
	buffRead <- waitingString
}
func readBuff(whattoexpect string, sshOut io.Reader, timeoutSeconds int) string {
	ch := make(chan string)
	go func(whattoexpect string, sshOut io.Reader) {
		buffRead := make(chan string)
		go readBuffForString(whattoexpect, sshOut, buffRead)
		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeoutSeconds) * time.Second):
			handleError(fmt.Errorf("%d", timeoutSeconds), true, fmt.Sprintf("Waiting for '%s' took longer than timeout: %d", whattoexpect, timeoutSeconds))
		}
	}(whattoexpect, sshOut)
	return <-ch
}
func writeBuff(command string, sshIn io.WriteCloser) (int, error) {
	returnCode, err := sshIn.Write([]byte(command + "\r"))
	return returnCode, err
}
func handleError(e error, fatal bool, customMessage ...string) {
	var errorMessage string
	if e != nil {
		if len(customMessage) > 0 {
			errorMessage = strings.Join(customMessage, " ")
		} else {
			errorMessage = "%s"
		}
		if fatal == true {
			log.Fatalf(errorMessage, e)
		} else {
			log.Print(errorMessage, e)
		}
	}
}

func loadCommands(fileName string) Commands {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	data := Commands{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}

	for _, cmd := range data.Commands {
		fmt.Println("Command: ", cmd)
	}

	return data
}

func main() {

	c := loadCommands("test/commands/cisco_ios.json")

	var results []string

	sshConfig := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.Password("Juniper"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.Config.Ciphers = append(sshConfig.Config.Ciphers, "aes128-cbc")
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	connection, err := ssh.Dial("tcp", "192.168.255.150:22", sshConfig)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	session, err := connection.NewSession()
	handleError(err, true, "Failed to create session: %s")
	sshOut, err := session.StdoutPipe()
	handleError(err, true, "Unable to setup stdin for session: %v")
	sshIn, err := session.StdinPipe()
	handleError(err, true, "Unable to setup stdout for session: %v")
	if err := session.RequestPty("xterm", 0, 200, modes); err != nil {
		session.Close()
		handleError(err, true, "request for pseudo terminal failed: %s")
	}

	if err := session.Shell(); err != nil {
		session.Close()
		handleError(err, true, "request for shell failed: %s")
	}
	readBuff("#", sshOut, 2)

	for _, cmd := range c.Commands {
		results = append(results, expecter(cmd, "#", 5, sshIn, sshOut))
	}

	fmt.Println(results)

	session.Close()
}
