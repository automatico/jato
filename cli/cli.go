package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/automatico/jato/command"
	"github.com/automatico/jato/device"
	"github.com/automatico/jato/user"
	"github.com/automatico/jato/utils"
	"golang.org/x/crypto/ssh/terminal"
)

// Params contain the result of CLI input
type Params struct {
	User     user.User
	Devices  device.Devices
	Commands command.Commands
	NoOp     bool
}

// CLI is the interface to the CLI application
func CLI() Params {
	userPtr := flag.String("u", os.Getenv("JATO_SSH_USER"), "Username to connect to devices with")
	askUserPassPtr := flag.Bool("p", false, "Ask for user password")
	devicesPtr := flag.String("d", "devices.json", "Devices inventory file")
	commandsPtr := flag.String("c", "commands.json", "Commands to run file")
	noOpPtr := flag.Bool("noop", false, "Dont execute job against devices")
	flag.Parse()

	// Collect CLI parameters
	p := Params{}

	// User
	p.User = user.User{}
	switch *userPtr != "" {
	case true:
		p.User.Username = *userPtr
	case false:
		fmt.Println("A username is required.")
		os.Exit(1)
	}

	// Password
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
	}
	p.User.Password = *userPass

	// Devices
	utils.FileStat(*devicesPtr)
	p.Devices = device.LoadDevices(*devicesPtr)

	// Commands
	utils.FileStat(*commandsPtr)
	p.Commands = command.LoadCommands(*commandsPtr)

	// No Op
	p.NoOp = *noOpPtr

	// CLI output
	t := utils.LoadTemplate("templates/cliRunner.templ")
	err = t.Execute(os.Stdout, p)
	if err != nil {
		panic(err)
	}

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
