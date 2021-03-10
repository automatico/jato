package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/automatico/jato/cli"
	"github.com/automatico/jato/device"
	"github.com/automatico/jato/output"
	"github.com/automatico/jato/telnet"
)

var telnetDevices []device.Device
var sshDevices []device.Device

func main() {

	cliParams := cli.CLI()

	for _, d := range cliParams.Devices.Devices {
		switch d.Connector {
		case "telnet":
			telnetDevices = append(telnetDevices, d)
		case "ssh":
			sshDevices = append(sshDevices, d)
		}
	}

	if cliParams.NoOp != true {
		//ssh.SSH(cliParams)
		results := telnet.Telnet(telnetDevices)
		t, err := template.New("results").Parse(output.CliResult)

		if err != nil {
			panic(err)
		}

		fmt.Println(output.JobResult)

		for _, r := range results.Results {
			err = t.Execute(os.Stdout, r)

			if err != nil {
				panic(err)
			}
		}

	}

}
