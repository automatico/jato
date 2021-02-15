package cli

import (
	"flag"
	"fmt"
	"os"
)

func CLI() {
	devicesPtr := flag.String("devices", "devices.json", "Devices inventory file")
	commandsPtr := flag.String("commands", "commands.json", "Commands to run file")
	usernamePtr := flag.String("username", os.Getenv("JATO_SSH_USER"), "Username to connect to devices")
	passwordPtr := flag.String("password", os.Getenv("JATO_SSH_PASS"), "Password for User")
	flag.Parse()

	fmt.Println("Devices: ", *devicesPtr)
	fmt.Println("Commands: ", *commandsPtr)
	fmt.Println("Username: ", *usernamePtr)
	fmt.Println("Password: ", *passwordPtr)
}
