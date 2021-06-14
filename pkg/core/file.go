package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/automatico/jato/internal/logger"
	"github.com/automatico/jato/internal/terminal"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/driver"
)

func LoadCommands(fileName string) data.Commands {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	commands := data.Commands{}

	err = json.Unmarshal([]byte(file), &commands)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	return commands
}

// Load a list of devices from a JSON file
func LoadDevices(fileName string) driver.Devices {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	devices := driver.Devices{}

	err = json.Unmarshal([]byte(file), &devices)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	return devices
}

func LoadVariables(fileName string) data.Variables {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

	variables := data.Variables{}

	err = json.Unmarshal([]byte(file), &variables)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	return variables
}

// writeStringToFile takes a string "s" and writes
// it to a writer "w"
func writeStringToFile(w *bufio.Writer, s string) {
	_, err := w.WriteString(s)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
}

// Write the output from commands run against
// devices to a plain text file
func WriteToFile(results []data.Result) {
	outdir := "output"
	for _, result := range results {

		CreateDeviceDir(fmt.Sprintf("%s/%s", outdir, result.Device))
		file, err := os.OpenFile(fmt.Sprintf("%s/%s/%d.raw", outdir, result.Device, result.Timestamp), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err))
		}

		writer := bufio.NewWriter(file)

		writeStringToFile(writer, fmt.Sprintf("! Device:    %s\n", result.Device))
		writeStringToFile(writer, fmt.Sprintf("! Timestamp: %d\n", result.Timestamp))
		writeStringToFile(writer, fmt.Sprintf("! OK:        %t\n", result.OK))
		writeStringToFile(writer, fmt.Sprintf("! Error:     %s\n", result.Error))

		for _, output := range result.CommandOutputs {

			writeStringToFile(writer, terminal.Banner(output.Command))
			writeStringToFile(writer, output.Output)
			writeStringToFile(writer, "\r\n")
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
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%d.json", outdir, result.Device, result.Timestamp), file, 0644)
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err))
		}
	}
}

// Create device directory if it does not
// already exist
func CreateDeviceDir(s string) {
	if _, err := os.Stat(s); os.IsNotExist(err) {
		err := os.MkdirAll(s, 0755)
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err))
		}
	}
}
