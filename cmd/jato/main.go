package main

import (
	"fmt"
	"html/template"
	"os"
	"sync"

	"github.com/automatico/jato/internal/templates"
	"github.com/automatico/jato/internal/terminal"
	"github.com/automatico/jato/pkg/jato"
)

var ciscoIOSDevices []jato.CiscoIOSDevice

func main() {

	// timeNow := time.Now().Unix()

	cliParams := jato.CLI()

	// Output data to feed into template
	data := map[string]interface{}{}
	data["banner"] = terminal.Banner("Job Parameters")
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
		d.Credentials.Username = cliParams.Credentials.Username
		d.Credentials.Password = cliParams.Credentials.Password
		d.Credentials.SuperPassword = cliParams.Credentials.SuperPassword
		d.Credentials.SSHKeyFile = cliParams.Credentials.SSHKeyFile
		if d.Vendor == "cisco" {
			if d.Platform == "ios" {
				cd := jato.NetToCiscoIOSDevice(d)
				ciscoIOSDevices = append(ciscoIOSDevices, cd)
			}
		}
	}

	if !cliParams.NoOp {

		results := []jato.Result{}

		var wg sync.WaitGroup
		ch := make(chan jato.Result)
		defer close(ch)

		wg.Add(len(ciscoIOSDevices))
		for _, dev := range ciscoIOSDevices {
			dev := dev // lock the host or the same host can run more than once
			dev.Init()
			if dev.Connector == "telnet" {
				go jato.RunWithTelnet(&dev, cliParams.Commands.Commands, ch, &wg)
			} else if dev.Connector == "ssh" {
				go jato.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
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

		jato.WriteToFile(results)
		jato.WriteToJSONFile(results)
	}

}
