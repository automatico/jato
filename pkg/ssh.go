package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type user struct {
	username string
	password string
}

type device struct {
	name     string
	vendor   string
	platform string
}

type commands struct {
	driver   string
	commands []string
}

// https://stackoverflow.com/a/63759067
// https://kukuruku.co/post/ssh-commands-execution-on-hundreds-of-servers-via-go/
// https://zaiste.net/posts/executing-commands-via-ssh-using-go/
func main() {

	commands := map[string][]string{
		"cisco": {
			"terminal length 0",
			"show version",
			"show ip interface brief",
			"show ip arp",
			"show cdp neighbors",
			"show running-config",
			"exit",
		},
		"juniper": {
			"set cli screen-length 0",
			"show interfaces terse",
			"show version",
			"show lldp neighbors",
			"exit",
		},
	}

	u := user{
		username: "admin",
		password: "Juniper",
	}

	devices := []map[string]string{
		{
			"name":     "192.168.255.150",
			"vendor":   "cisco",
			"platform": "ios",
		},
		{
			"name":     "192.168.255.154",
			"vendor":   "cisco",
			"platform": "ios",
		},
		{
			"name":     "192.168.255.155",
			"vendor":   "cisco",
			"platform": "ios",
		},
		{
			"name":     "192.168.255.156",
			"vendor":   "cisco",
			"platform": "ios",
		},
		{
			"name":     "192.168.255.157",
			"vendor":   "cisco",
			"platform": "ios",
		},
	}

	sshconfig := InsecureClientConfig(u.username, u.password)

	results := make(chan map[string]map[string]string)
	timeout := time.After(10 * time.Second)

	for _, device := range devices {
		go func(device string) {
			results <- ExecCommands(device, commands["cisco"], sshconfig)
			// result, _ := ExecCommands(device, listCMDs, sshconfig)
			// printResult(result)
		}(device["name"])
	}
	for i := 0; i < len(devices); i++ {
		select {
		case res := <-results:

			writeToJSONFile(res)
			// writeToFile(res)
			// printResult(res)

		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}

// ExecCommands ...
// func ExecCommands(ipAddr string, commands []string, sshconfig *ssh.ClientConfig) ([]string, error) {
func ExecCommands(ipAddr string, commands []string, sshconfig *ssh.ClientConfig) map[string]map[string]string {

	// Gets IP, credentials and config/commands, SSH Config (Timeout, Ciphers, ...) and returns
	// output of the device as "string" and an error. If error == nil, means program was able to SSH with no issue

	// Creating outerr as Output Error.
	// outerr := errors.New("nil")
	// outerr = nil
	// fmt.Println(outerr)

	// Creating Output as String
	var outputStr []string
	var strTmp string

	results := make(map[string]string)

	// Dial to the remote-host
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ipAddr), sshconfig)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create sesssion
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// StdinPipe() returns a pipe that will be connected to the remote command's standard input when the command starts.
	// StdoutPipe() returns a pipe that will be connected to the remote command's standard output when the command starts.
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start remote shell
	err = session.Shell()
	if err != nil {
		log.Fatal(err)
	}

	stdinLines := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdinLines <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Printf("scanner failed: %v", err)
		}
		close(stdinLines)
	}()

	// Send the commands to the remotehost one by one.
	for i, cmd := range commands {

		command := underscorize(cmd)

		_, err := stdin.Write([]byte(fmt.Sprintf("%s\n", cmd)))
		if err != nil {
			log.Fatal(err)
		}
		if i == len(commands)-1 {
			_ = stdin.Close() // send eof
		}

		// wait for command to complete
		// we'll assume the moment we've gone 1 secs w/o any output that our command is done
		timer := time.NewTimer(0)
	InputLoop:
		for {
			timer.Reset(time.Millisecond * 1000)
			select {
			case line, ok := <-stdinLines:
				if !ok {
					log.Println("Finished processing")
					break InputLoop
				}
				strTmp += line
				strTmp += "\n"
			case <-timer.C:
				break InputLoop
			}
		}
		outputStr = append(outputStr, strTmp)
		results[command] = strTmp
		//log.Printf("Finished processing %v\n", cmd)
		strTmp = ""
	}
	// Wait for session to finish
	err = session.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// return outputStr, outerr
	result := map[string]map[string]string{
		ipAddr: results,
	}
	return result
}

// InsecureClientConfig ...
func InsecureClientConfig(userStr, passStr string) *ssh.ClientConfig {

	SSHconfig := &ssh.ClientConfig{
		User:    userStr,
		Timeout: 5 * time.Second,
		Auth:    []ssh.AuthMethod{ssh.Password(passStr)},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// Config: ssh.Config{
		//     Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-cbc", "aes192-cbc",
		//         "aes256-cbc", "3des-cbc", "des-cbc"},
		//     KeyExchanges: []string{"diffie-hellman-group1-sha1",
		//         "diffie-hellman-group-exchange-sha1",
		//         "diffie-hellman-group14-sha1"},
		// },
	}
	return SSHconfig
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
