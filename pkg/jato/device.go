package jato

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Device represents a managed device
type Device struct {
	IP        string `json:"ip"`
	Name      string `json:"name"`
	Vendor    string `json:"vendor"`
	Platform  string `json:"platform"`
	Connector string `json:"connector"`
}

// Devices holds a collection of Device structs
type Devices struct {
	Devices []Device `json:"devices"`
}

// Load a list of devices from a JSON file
func LoadDevices(fileName string) Devices {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	data := Devices{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (d Device) Auth(c UserCredentials) {}
