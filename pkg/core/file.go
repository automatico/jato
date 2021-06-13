package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/automatico/jato/internal/terminal"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/driver"
)

func LoadCommands(fileName string) data.Commands {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	commands := data.Commands{}

	err = json.Unmarshal([]byte(file), &commands)
	if err != nil {
		log.Fatal(err)
	}

	return commands
}

// Load a list of devices from a JSON file
func LoadDevices(fileName string) driver.Devices {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	devices := driver.Devices{}

	err = json.Unmarshal([]byte(file), &devices)
	if err != nil {
		log.Fatal(err)
	}

	return devices
}

func LoadVariables(fileName string) data.Variables {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	variables := data.Variables{}

	err = json.Unmarshal([]byte(file), &variables)
	if err != nil {
		log.Fatal(err)
	}
	return variables
}

// Write the output from commands run against
// devices to a plain text file
func WriteToFile(results []data.Result) {
	outdir := "output"
	for _, result := range results {
		CreateDeviceDir(fmt.Sprintf("%s/%s", outdir, result.Device))
		file, err := os.OpenFile(fmt.Sprintf("%s/%s/%d.raw", outdir, result.Device, result.Timestamp), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writer := bufio.NewWriter(file)

		_, err = writer.WriteString(fmt.Sprintf("! Device:    %s\n", result.Device))
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
		_, err = writer.WriteString(fmt.Sprintf("! Timestamp: %d\n", result.Timestamp))
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
		_, err = writer.WriteString(fmt.Sprintf("! OK:        %t\n", result.OK))
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
		_, err = writer.WriteString(fmt.Sprintf("! Error:     %s\n", result.Error))
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}

		for _, output := range result.CommandOutputs {
			_, err = writer.WriteString(terminal.Banner(output.Command))
			if err != nil {
				log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
			}
			_, err = writer.WriteString(output.Output)
			if err != nil {
				log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
			}
			_, err = writer.WriteString("\r\n")
			if err != nil {
				log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
			}
		}
		writer.Flush()
	}
}

// Write the output from commands run against
// devices to a json file
func WriteToJSONFile(results []data.Result) {
	outdir := "output"
	for _, result := range results {
		CreateDeviceDir(fmt.Sprintf("%s/%s", outdir, result.Device))
		file, _ := json.MarshalIndent(result, "", " ")
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%s/%d.json", outdir, result.Device, result.Timestamp), file, 0644)
	}
}

// Create device directory if it does not
// already exist
func CreateDeviceDir(s string) {
	if _, err := os.Stat(s); os.IsNotExist(err) {
		err := os.MkdirAll(s, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
