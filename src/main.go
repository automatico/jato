package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"os"
	"bufio"
	"regexp"

	"golang.org/x/crypto/ssh"
	// Uncomment to store output in variable
	//"bytes"
)

type user struct {
	username string
	password string
}

type device struct {
	name string
	vendor string
	platform string
}

func executeCmd(hostname string, cmds []string, config *ssh.ClientConfig) *[]string {
	// Need pseudo terminal if we want to have an SSH session
	// similar to what you have when you use a SSH client
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	conn, err := ssh.Dial("tcp", hostname, config)
	if err != nil {
		log.Println(err)
		return &[]string{1: hostname}
	}
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	// You can use session.Run() here but that only works
	// if you need a run a single command or you commands
	// are independent of each other.
	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}
	stdoutBuf, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("request for stdout pipe failed: %s", err)
	}
	stdinBuf, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("request for stdin pipe failed: %s", err)
	}
	err = session.Shell()
	if err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}

	for _, cmd := range cmds {
		stdinBuf.Write([]byte(cmd + "\n"))
	}
	res := make([]string, 0)
	return readStdoutBuf(stdoutBuf, &res, hostname)
}

func readStdoutBuf(stdBuf io.Reader, res *[]string, hostname string) *[]string {
	stdoutBuf := make([]byte, 1000000)
	time.Sleep(time.Millisecond * 100)
	byteCount, err := stdBuf.Read(stdoutBuf)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Bytes received: ", byteCount)
	s := string(stdoutBuf[:byteCount])
	lines := strings.Split(s, "\n")
	
	// writeToFile(lines)
	printResult(lines)
	
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()

	// fmt.Println("Here")
	// fmt.Println(strings.Contains(strings.TrimSpace(lines[len(lines)-1]), "iosv#"))
	// fmt.Println(lines[len(lines)-1])

	re := regexp.MustCompile(`^.*#$`)
	if strings.TrimSpace(lines[len(lines)-1]) != re.String() {
		fmt.Println()
		fmt.Println(lines[len(lines)-1])
		*res = append(*res, lines...)
		readStdoutBuf(stdBuf, res, hostname)
		// return res
	}
	fmt.Println("end reached")
	*res = append(*res, lines...)
	return res
}

func main() {

	commands := map[string][]string{
		"cisco": {
			"term len 0",
			"show ip int brie",
			"show version",
			"show run",
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

	ciscoDevice := device{
		name: "192.168.255.150",
		vendor: "cisco",
		platform: "ios",
	}

	var outStrings []string

	// juniperDevice := device{
	// 	name: "192.168.255.151",
	// 	vendor: "juniper",
	// 	platform: "junios",
	// }
// 
	// aristaDevice := device{
	// 	name: "192.168.255.152",
	// 	vendor: "arista",
	// 	platform: "eos",
	// }

	port := "22"

	// SSH client config
	config := &ssh.ClientConfig{
		User: u.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(u.password),
		},
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Delete temp file
	deleteFile("temp.txt")

	// Connect to host
	host := fmt.Sprintf("%s:%s", ciscoDevice.name, port)

	results := make(chan *[]string, 100)
	go func(hostname string) {
		results <- executeCmd(hostname, commands["cisco"], config)
	}(host)

	res := <-results
	outStrings = append(outStrings, *res...)

	printResult(outStrings)
}

func printResult(result []string) {
	for _, line := range result {
		fmt.Println(line)
	}
}

func writeToFile(lines []string) {
	file, err := os.OpenFile("temp.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
			log.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
					log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
			}
	}
	writer.Flush()

}

func deleteFile(filename string) {
	err := os.Remove(filename)  // remove a single file
	if err != nil {
	  fmt.Println(err)
	}
}