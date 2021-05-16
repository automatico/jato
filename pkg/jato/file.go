package jato

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func LoadCommands(fileName string) CommandExpect {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	data := CommandExpect{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
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

// Write the output from commands run against
// devices to a plain text file
func WriteToFile(timestamp int64, results []Result) {
	outdir := "data"
	for _, result := range results {
		CreateDeviceDir(fmt.Sprintf("%s/%s", outdir, result.Device))
		file, err := os.OpenFile(fmt.Sprintf("%s/%s/%d.raw", outdir, result.Device, timestamp), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writer := bufio.NewWriter(file)
		for _, output := range result.CommandOutputs {
			_, err := writer.WriteString(output.Output)
			if err != nil {
				log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
			}
		}
		writer.Flush()
		fmt.Printf("Saved RAW output: %s\n", result.Device)
	}
}

// Write the output from commands run against
// devices to a json file
func WriteToJSONFile(timestamp int64, results []Result) {
	// fmt.Println("########################")
	// fmt.Println(results)
	// fmt.Println("########################")

	outdir := "data"
	for _, result := range results {
		CreateDeviceDir(fmt.Sprintf("%s/%s", outdir, result.Device))
		file, _ := json.MarshalIndent(result, "", " ")
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%s/%d.json", outdir, result.Device, timestamp), file, 0644)
		fmt.Printf("Saved JSON output: %s\n", result.Device)
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
