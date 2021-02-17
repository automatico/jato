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
