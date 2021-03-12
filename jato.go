package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/automatico/jato/cli"
	"github.com/automatico/jato/connector"
	"github.com/automatico/jato/device"
	"github.com/automatico/jato/output"
	"github.com/automatico/jato/telnet"
	"github.com/automatico/jato/templates"
)

var telnetDevices []device.Device
var sshDevices []device.Device

func main() {

	cliParams := cli.CLI()

	jt := connector.Jato{
		UserCredentials: cliParams.Credentials,
		Devices:         cliParams.Devices,
		CommandExpect:   cliParams.Commands,
	}

	// Output data to feed into template
	data := map[string]interface{}{}
	data["divider"] = output.Divider("Job Parameters")
	data["params"] = cliParams

	// CLI output
	t, err := template.New("output").Parse(templates.CliRunner)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, data)

	if err != nil {
		panic(err)
	}

	for _, d := range cliParams.Devices.Devices {
		switch d.Connector {
		case "telnet":
			telnetDevices = append(telnetDevices, d)
		case "ssh":
			sshDevices = append(sshDevices, d)
		}
	}
	jt.Devices.Devices = telnetDevices
	fmt.Println(telnetDevices)

	if !cliParams.NoOp {
		// ssh.SSH(cliParams)
		results := telnet.Telnet(jt)
		t, err := template.New("results").Parse(templates.CliResult)

		if err != nil {
			panic(err)
		}

		fmt.Print(output.Divider("Job Results"))

		for _, r := range results.Results {
			err = t.Execute(os.Stdout, r)

			if err != nil {
				panic(err)
			}
		}
	}

}
