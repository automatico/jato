package command

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Commands to run against a device
type Commands struct {
	Commands []string `json:"commands"`
}

// CommandExpect is a command to run
// with the expected string to match
type CommandExpect struct {
	Command string
	Expect  string
}

// LoadCommands is used to load a hash of commands
// from a json file in the following format.
// {
// 	"commands": [
// 	  "terminal length 0",
// 	  "show version",
// 	  "show ip interface brief",
// 	  "show ip arp",
// 	  "show cdp neighbors",
// 	  "show running-config",
// 	  "exit"
// 	]
//}
func LoadCommands(fileName string) Commands {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	data := Commands{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
