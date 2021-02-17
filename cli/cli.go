package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/automatico/jato/command"
	"github.com/automatico/jato/device"
	"github.com/automatico/jato/user"
	"golang.org/x/crypto/ssh/terminal"
)

type CLIParams struct {
	User     user.User
	Devices  device.Devices
	Commands command.Commands
}

func CLI() CLIParams {
	userPtr := flag.String("u", os.Getenv("JATO_SSH_USER"), "Username to connect to devices with")
	askUserPassPtr := flag.Bool("p", false, "Ask for user password")
	devicesPtr := flag.String("d", "devices.json", "Devices inventory file")
	commandsPtr := flag.String("c", "commands.json", "Commands to run file")
	flag.Parse()

	p := CLIParams{}
	p.User = user.User{}

	if *userPtr == "" {
		fmt.Println("A username is required.")
		os.Exit(1)
	} else {
		p.User.Username = *userPtr
	}

	userPass := new(string)
	var err error
	if *askUserPassPtr == true {
		*userPass, err = promptSecret("Enter user password?")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if *askUserPassPtr == false {
		*userPass = os.Getenv("JATO_SSH_PASS")
		if *userPass == "" {
			fmt.Println("A password is required.")
			os.Exit(1)
		}
	} else {
		p.User.Password = *userPass
	}

	fileStat(*devicesPtr)
	p.Devices = device.LoadDevices(*devicesPtr)

	fileStat(*commandsPtr)
	p.Commands = command.LoadCommands(*commandsPtr)

	fmt.Println("Username: ", *userPtr)
	fmt.Println("Password: ", "********")
	fmt.Println("Devices:  ", *devicesPtr)
	fmt.Println("Commands: ", *commandsPtr)

	return p
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

func fileStat(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Filename: '%s' does not exist or is not readable.", filename)
		os.Exit(1)
	}
}
