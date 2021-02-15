package cli

import (
	"flag"
	"fmt"
)

func CLI() {
	devicesPtr := flag.String("devices", "devices.json", "Devices inventory file")
	commandsPtr := flag.String("commands", "commands.json", "Commands to run file")

	flag.Parse()

	fmt.Println("Devices: ", *devicesPtr)
	fmt.Println("Commands: ", *commandsPtr)
}
