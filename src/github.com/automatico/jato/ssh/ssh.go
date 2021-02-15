package ssh

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
	IP        string `json:"ip"`
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

func readBuffForString(expect string, sshOut io.Reader, buffRead chan<- string) {
	buf := make([]byte, 1024)
	n, err := sshOut.Read(buf) //this reads the ssh terminal
	tmpStr := ""
	if err == nil {
		tmpStr = string(buf[:n])
	}
	for (err == nil) && (!strings.Contains(tmpStr, expect)) {
		n, err = sshOut.Read(buf)
		tmpStr += string(buf[:n])
		// fmt.Println(tmpStr) //uncommenting this might help you debug if you are coming into errors with timeouts when correct details entered

	}
	buffRead <- tmpStr
}
func readBuff(expect string, sshOut io.Reader, timeout int) string {
	ch := make(chan string)
	go func(expect string, sshOut io.Reader) {
		buffRead := make(chan string)
		go readBuffForString(expect, sshOut, buffRead)
		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeout) * time.Second):
			handleError(fmt.Errorf("%d", timeout), true, fmt.Sprintf("Waiting for '%s' took longer than timeout: %d", expect, timeout))
		}
	}(expect, sshOut)
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
func writeToFile(timestamp int64, results map[string]map[string]string) {
	outdir := "data"
	for k, v := range results {
		createDeviceDir(fmt.Sprintf("%s/%s", outdir, k))
		file, err := os.OpenFile(fmt.Sprintf("%s/%s/%d.raw", outdir, k, timestamp), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		fmt.Printf("Saved RAW output: %s\n", k)
	}
}

// Write the output from commands run against
// devices to a json file
func writeToJSONFile(timestamp int64, results map[string]map[string]string) {
	outdir := "data"
	for k, v := range results {
		createDeviceDir(fmt.Sprintf("%s/%s", outdir, k))
		file, _ := json.MarshalIndent(v, "", " ")
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%s/%d.json", outdir, k, timestamp), file, 0644)
		fmt.Printf("Saved JSON output: %s\n", k)
	}
}

// Converts a string to an underscore string
// replacing spaces and dashes with underscores
func underscorer(s string) string {
	re := strings.NewReplacer(" ", "_", "-", "_")
	return re.Replace(s)
}

// Create device directory if it does not
// already exist
func createDeviceDir(s string) {
	if _, err := os.Stat(s); os.IsNotExist(err) {
		err := os.MkdirAll(s, 0755)
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
		// Make this an option
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Make this an option
	sshConfig.Config.Ciphers = append(sshConfig.Config.Ciphers, "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-cbc", "aes192-cbc", "aes256-cbc", "3des-cbc", "des-cbc")
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", device.IP, SSHPort), sshConfig)
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

		results[underscorer(cmd)] = expecter(cmd, "#", 5, sshIn, sshOut)
	}
	session.Close()
	res := map[string]map[string]string{
		device.Name: results,
	}
	return res
}

func SSH() {

	t := time.Now().Unix()
	user := User{
		username: os.Getenv("JATO_SSH_USER"),
		password: os.Getenv("JATO_SSH_PASS"),
	}
	commands := loadCommands("commands.json")
	devices := loadDevices("devices.json")

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
			writeToJSONFile(t, res)
			writeToFile(t, res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}

}
