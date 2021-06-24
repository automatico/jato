package core

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/automatico/jato/internal/logger"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/driver"
	"golang.org/x/term"
)

const version = "2021.06.13"

// Params contain the result of CLI input
type Params struct {
	Credentials data.Credentials
	Devices     driver.NetDevices
	Commands    data.Commands
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

	vars := data.Variables{}

	if *versionPtr {
		fmt.Printf("jato version: %s\n", version)
		os.Exit(0)
	}

	// Used to collect CLI parameters
	params := Params{}

	userCreds := data.GetCredentials(vars.Credentials)

	// User
	params.Credentials = userCreds

	if *userPtr != "" {
		params.Credentials.Username = *userPtr
	} else if params.Credentials.Username == "" {
		logger.Fatal("a username is required")
	}

	// Password
	userPass := new(string)
	var err error
	if *askUserPassPtr {
		*userPass, err = promptSecret("Enter user password:")
		params.Credentials.Password = *userPass
		if err != nil {
			logger.Fatal(err)
		}
	} else if !*askUserPassPtr {
		if userCreds.Password == "" {
			logger.Fatal("a password is required")
		}
	}

	// Devices
	if err := FileStat(*devicesPtr); err != nil {
		logger.Fatalf("device file does not exist: %v", *devicesPtr)
	}
	params.Devices = LoadDevices(*devicesPtr)

	// Commands
	if err := FileStat(*commandsPtr); err != nil {
		logger.Fatalf("command file does not exist: %v", *commandsPtr)
	}
	params.Commands = LoadCommands(*commandsPtr)

	// No Op
	params.NoOp = *noOpPtr

	return params
}

// promptSecret prompts user for an input that is not echo-ed on terminal.
func promptSecret(question string) (string, error) {
	fmt.Printf(question + "\n=> ")

	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	pass := string(bytepw)
	return pass, nil
}
