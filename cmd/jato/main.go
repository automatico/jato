package main

import (
	"fmt"
	"html/template"
	"os"
	"sync"

	"github.com/automatico/jato/internal/templates"
	"github.com/automatico/jato/internal/terminal"
	"github.com/automatico/jato/pkg/core"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/driver"
	"github.com/automatico/jato/pkg/network"
)

var ciscoIOSDevices []driver.CiscoIOSDevice

func main() {

	// timeNow := time.Now().Unix()

	cliParams := core.CLI()

	// Output data to feed into template
	templateData := map[string]interface{}{}
	templateData["banner"] = terminal.Banner("Job Parameters")
	templateData["params"] = cliParams

	// CLI output
	t, err := template.New("output").Parse(templates.CliRunner)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, templateData)

	if err != nil {
		panic(err)
	}

	for _, d := range cliParams.Devices.Devices {
		d.Credentials.Username = cliParams.Credentials.Username
		d.Credentials.Password = cliParams.Credentials.Password
		d.Credentials.SuperPassword = cliParams.Credentials.SuperPassword
		d.Credentials.SSHKeyFile = cliParams.Credentials.SSHKeyFile
		if d.Vendor == "cisco" {
			if d.Platform == "ios" {
				cd := driver.NewCiscoIOSDevice(d)
				ciscoIOSDevices = append(ciscoIOSDevices, cd)
			}
		}
	}

	if !cliParams.NoOp {

		results := []data.Result{}

		var wg sync.WaitGroup
		ch := make(chan data.Result)
		defer close(ch)

		wg.Add(len(ciscoIOSDevices))
		for _, dev := range ciscoIOSDevices {
			dev := dev // lock the host or the same host can run more than once
			dev.Init()
			if dev.Connector == "telnet" {
				go network.RunWithTelnet(&dev, cliParams.Commands.Commands, ch, &wg)
			} else if dev.Connector == "ssh" {
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		devTotal := len(ciscoIOSDevices)
		for i := 0; i < devTotal; i++ {
			results = append(results, <-ch)
		}

		wg.Wait()

		t, err := template.New("results").Parse(templates.CliResult)
		if err != nil {
			panic(err)
		}

		fmt.Print(terminal.Banner("Job Results"))

		for _, r := range results {
			err = t.Execute(os.Stdout, r)

			if err != nil {
				panic(err)
			}
		}

		core.WriteToFile(results)
		core.WriteToJSONFile(results)
	}

}
