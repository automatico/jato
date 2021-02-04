package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
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

func main() {

	u := user{
		username: "admin",
		password: "Juniper",
	}

	ciscoDevice := device{
		name: "192.168.255.150",
		vendor: "cisco",
		platform: "ios",
	}

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

	// modes := ssh.TerminalModes{
	// 	ssh.ECHO:          0,     // disable echoing
	// 	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	// 	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	// }

	// Connect to host
	host := fmt.Sprintf("%s:%s", ciscoDevice.name, port)

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create sesssion
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer sess.Close()

	// StdinPipe for commands
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment to store output in variable
	//var b bytes.Buffer
	//sess.Stdout = b
	//sess.Stderr = b

	// Enable system stdout
	// Comment these if you uncomment to store in variable
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// Start remote shell
	err = sess.Shell()
	if err != nil {
		log.Fatal(err)
	}

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
	for _, cmd := range commands["cisco"] {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Wait for sess to finish
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment to store in variable
	//fmt.Println(b.String())

}
