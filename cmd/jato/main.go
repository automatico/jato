package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/automatico/jato/internal"
	"github.com/automatico/jato/internal/templates"
	"github.com/automatico/jato/pkg/jato"
)

type Jato struct {
	jato.Credentials
	jato.Devices
	jato.CommandExpect
}

var telnetDevices []jato.Device
var sshDevices []jato.Device

func main() {

	cliParams := jato.CLI()

	jt := Jato{
		Credentials:   cliParams.Credentials,
		Devices:       cliParams.Devices,
		CommandExpect: cliParams.Commands,
	}

	// Output data to feed into template
	data := map[string]interface{}{}
	data["divider"] = internal.Divider("Job Parameters")
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

	if !cliParams.NoOp {
		// ssh.SSH(cliParams)
		results := jato.Telnet(jt)
		t, err := template.New("results").Parse(templates.CliResult)

		if err != nil {
			panic(err)
		}

		fmt.Print(internal.Divider("Job Results"))

		for _, r := range results.Results {
			err = t.Execute(os.Stdout, r)

			if err != nil {
				panic(err)
			}
		}
	}

}
