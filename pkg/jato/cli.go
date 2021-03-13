package jato

import (
	"flag"
	"fmt"
	"os"

	"github.com/automatico/jato/internal"
	"golang.org/x/crypto/ssh/terminal"
)

const version = "2021.02.02"

// Params contain the result of CLI input
type Params struct {
	Credentials UserCredentials
	Devices     Devices
	Commands    CommandExpect
	NoOp        bool
}

// CLI is the interface to the CLI application
func CLI() Params {
	userPtr := flag.String("u", os.Getenv("JATO_SSH_USER"), "Username to connect to devices with")
	askUserPassPtr := flag.Bool("a", false, "Ask for user password")
	devicesPtr := flag.String("d", "devices.json", "Devices inventory file")
	commandsPtr := flag.String("c", "commands.json", "Commands to run file")
	noOpPtr := flag.Bool("noop", false, "Don't execute job against devices")
	versionPtr := flag.Bool("v", false, "Jato version")
	flag.Parse()

	if *versionPtr {
		fmt.Printf("Jato version: %s\n", version)
		os.Exit(0)
	}

	// Used to collect CLI parameters
	params := Params{}

	userCreds := UserCredentials{}.Load()

	// User
	params.Credentials = userCreds

	if *userPtr != "" {
		params.Credentials.Username = *userPtr
	} else if params.Credentials.Username == "" {
		fmt.Println("A username is required.")
		os.Exit(1)
	}

	// Password
	userPass := new(string)
	var err error
	if *askUserPassPtr {
		*userPass, err = promptSecret("Enter user password:")
		params.Credentials.Password = *userPass
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if !*askUserPassPtr {
		if userCreds.Password == "" {
			fmt.Println("A password is required.")
			os.Exit(1)
		}
	}

	// Devices
	internal.FileStat(*devicesPtr)
	params.Devices = LoadDevices(*devicesPtr)

	// Commands
	internal.FileStat(*commandsPtr)
	params.Commands = LoadCommands(*commandsPtr)

	// No Op
	params.NoOp = *noOpPtr

	return params
}

// promptSecret prompts user for an input that is not echo-ed on terminal.
func promptSecret(question string) (string, error) {
	fmt.Printf(question + "\n> ")

	raw, err := terminal.MakeRaw(0)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(0, raw)

	var (
		prompt string
		answer string
	)

	term := terminal.NewTerminal(os.Stdin, prompt)
	for {
		char, err := term.ReadPassword(prompt)
		if err != nil {
			return "", err
		}
		answer += char

		if char == "" || char == answer {
			return answer, nil
		}
	}
}
