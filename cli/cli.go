package cli

import (
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/automatico/jato/command"
	"github.com/automatico/jato/credentials"
	"github.com/automatico/jato/device"
	"github.com/automatico/jato/output"
	"github.com/automatico/jato/utils"
	"golang.org/x/crypto/ssh/terminal"
)

const version = "2021.02.02"

// Params contain the result of CLI input
type Params struct {
	Credentials credentials.UserCredentials
	Devices     device.Devices
	Commands    command.Commands
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

	if *versionPtr == true {
		fmt.Printf("Jato version: %s\n", version)
		os.Exit(0)
	}

	// Used to collect CLI parameters
	params := Params{}

	userCreds := credentials.UserCredentials{}.Load()

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
	if *askUserPassPtr == true {
		*userPass, err = promptSecret("Enter user password:")
		params.Credentials.Password = *userPass
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if *askUserPassPtr == false {
		if userCreds.Password == "" {
			fmt.Println("A password is required.")
			os.Exit(1)
		}
	}

	// Devices
	utils.FileStat(*devicesPtr)
	params.Devices = device.LoadDevices(*devicesPtr)

	// Commands
	utils.FileStat(*commandsPtr)
	params.Commands = command.LoadCommands(*commandsPtr)

	// No Op
	params.NoOp = *noOpPtr

	// Output data to feed into template
	data := map[string]interface{}{}
	data["divider"] = output.Divider("Job Parameters")
	data["params"] = params

	// CLI output
	t, err := template.New("output").Parse(output.CliRunner)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, data)

	if err != nil {
		panic(err)
	}

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
