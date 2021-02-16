package cli

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func CLI() {
	userPtr := flag.String("u", os.Getenv("JATO_SSH_USER"), "Username to connect to devices with")
	askUserPass := flag.Bool("p", false, "Ask for user password")
	devicesPtr := flag.String("d", "devices.json", "Devices inventory file")
	commandsPtr := flag.String("c", "commands.json", "Commands to run file")
	flag.Parse()

	if *userPtr == "" {
		fmt.Println("A username is required.")
		os.Exit(1)
	}

	userPass := new(string)
	var err error
	if *askUserPass == true {
		*userPass, err = promptSecret("Enter user password?")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if *askUserPass == false {
		*userPass = os.Getenv("JATO_SSH_PASS")
		if *userPass == "" {
			fmt.Println("A password is required.")
			os.Exit(1)
		}
	}

	fileStat(*devicesPtr)
	fileStat(*commandsPtr)

	fmt.Println("Username: ", *userPtr)
	fmt.Println("Password: ", "********")
	fmt.Println("Devices:  ", *devicesPtr)
	fmt.Println("Commands: ", *commandsPtr)

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
