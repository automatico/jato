package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHPort ...
const (
	SSHPort = 22
)

// User represents a users credentials
type User struct {
	username string
	password string
}

// Device represents a managed device
type Device struct {
	Name      string `json:"name"`
	Vendor    string `json:"vendor"`
	Platform  string `json:"platform"`
	Connector string `json:"connector"`
}

// Devices holds a collection of Device structs
type Devices struct {
	Device []Device `json:"devices"`
}

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
	buf := make([]byte, 1024)
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

	return data
}

// Print the output from commands run against
// devices to stdout
func printResult(result map[string]map[string]string) {
	for k, v := range result {
		fmt.Println(fmt.Sprintf("hostname: %s", k))
		fmt.Println(v)
	}
}

// Write the output from commands run against
// devices to a plain text file
func writeToFile(results map[string]map[string]string) {
	t := time.Now().Unix()
	for k, v := range results {
		createDeviceDir(k)
		file, err := os.OpenFile(fmt.Sprintf("%s/%d.raw", k, t), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writer := bufio.NewWriter(file)
		for _, output := range v {
			_, err := writer.WriteString(output)
			if err != nil {
				log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
			}
		}
		writer.Flush()
	}
}

// Write the output from commands run against
// devices to a json file
func writeToJSONFile(results map[string]map[string]string) {
	t := time.Now().Unix()
	for k, v := range results {
		createDeviceDir(k)
		file, _ := json.MarshalIndent(v, "", " ")
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%d.json", k, t), file, 0644)
	}
}

// Converts a string to an underscore string
// replacing spaces and dashes with underscores
func underscorize(s string) string {
	re := strings.NewReplacer(" ", "_", "-", "_")
	return re.Replace(s)
}

// Create device directory if it does not
// already exist
func createDeviceDir(s string) {
	if _, err := os.Stat(s); os.IsNotExist(err) {
		err := os.Mkdir(s, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Load a list of devices from a JSON file
func loadDevices(fileName string) Devices {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	data := Devices{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func runner(user User, device Device, commands Commands) map[string]map[string]string {

	results := make(map[string]string)

	sshConfig := &ssh.ClientConfig{
		User: user.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(user.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.Config.Ciphers = append(sshConfig.Config.Ciphers, "aes128-cbc")
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", device.Name, SSHPort), sshConfig)
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

	for _, cmd := range commands.Commands {

		results[underscorize(cmd)] = expecter(cmd, "#", 5, sshIn, sshOut)
	}
	session.Close()
	res := map[string]map[string]string{
		device.Name: results,
	}
	return res
}

func main() {

	user := User{
		username: "admin",
		password: "Juniper",
	}
	commands := loadCommands("test/commands/cisco_ios.json")
	devices := loadDevices("test/devices/cisco.json")

	results := make(chan map[string]map[string]string)
	timeout := time.After(10 * time.Second)

	for _, device := range devices.Device {
		go func(user User, device Device, commands Commands) {
			results <- runner(user, device, commands)
		}(user, device, commands)
	}

	for range devices.Device {
		select {
		case res := <-results:
			writeToJSONFile(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}

}
