package main

import (
	"fmt"
	"html/template"
	"os"
	"sync"

	"github.com/automatico/jato/internal/logger"
	"github.com/automatico/jato/internal/templates"
	"github.com/automatico/jato/internal/terminal"
	"github.com/automatico/jato/pkg/core"
	"github.com/automatico/jato/pkg/data"
	"github.com/automatico/jato/pkg/driver"
	"github.com/automatico/jato/pkg/network"
)

var ciscoIOSDevices []driver.CiscoIOSDevice
var aristaEOSDevices []driver.AristaEOSDevice

func main() {

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

		vendorPlatform := fmt.Sprintf("%s_%s", d.Vendor, d.Platform)
		switch vendorPlatform {
		case "cisco_ios":
			cd := driver.NewCiscoIOSDevice(d)
			ciscoIOSDevices = append(ciscoIOSDevices, cd)
		case "arista_eos":
			ad := driver.NewAristaEOSDevice(d)
			aristaEOSDevices = append(aristaEOSDevices, ad)
		default:
			logger.Warning(fmt.Sprintf("device: %s with vendor: %s and platform: %s not supported", d.Name, d.Vendor, d.Platform))
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
			switch dev.Connector {
			case "telnet":
				go network.RunWithTelnet(&dev, cliParams.Commands.Commands, ch, &wg)
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		wg.Add(len(aristaEOSDevices))
		for _, dev := range aristaEOSDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		devTotal := len(ciscoIOSDevices) + len(aristaEOSDevices)
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
