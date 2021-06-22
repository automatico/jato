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
)

var allDevices []driver.NetDevice

func main() {

	cliParams := core.CLI()

	// Output data to feed into template
	templateData := map[string]interface{}{}
	templateData["banner"] = terminal.Banner("Job Parameters")
	templateData["params"] = cliParams

	// CLI output
	t, err := template.New("output").Parse(templates.CliRunner)
	if err != nil {
		logger.Fatal(err)
	}

	if err := t.Execute(os.Stdout, templateData); err != nil {
		logger.Fatal(err)
	}

	for _, d := range cliParams.Devices.Devices {

		creds := d.Variables.Credentials
		if creds != "" {
			d.Credentials = data.GetCredentials(creds)
		} else {
			d.Credentials.Username = cliParams.Credentials.Username
			d.Credentials.Password = cliParams.Credentials.Password
			d.Credentials.SuperPassword = cliParams.Credentials.SuperPassword
			d.Credentials.SSHKeyFile = cliParams.Credentials.SSHKeyFile
		}

		vendorPlatform := fmt.Sprintf("%s_%s", d.Vendor, d.Platform)
		switch vendorPlatform {
		case "arista_eos":
			nd := driver.NewAristaEOSDevice(d)
			allDevices = append(allDevices, nd)
		case "aruba_aoscx":
			nd := driver.NewArubaAOSCXDevice(d)
			allDevices = append(allDevices, nd)
		case "cisco_aireos":
			nd := driver.NewCiscoAireOSDevice(d)
			allDevices = append(allDevices, nd)
		case "cisco_asa":
			nd := driver.NewCiscoASADevice(d)
			allDevices = append(allDevices, nd)
		case "cisco_ios":
			nd := driver.NewCiscoIOSDevice(d)
			allDevices = append(allDevices, nd)
		case "cisco_iosxr":
			nd := driver.NewCiscoIOSXRDevice(d)
			allDevices = append(allDevices, nd)
		case "cisco_nxos":
			nd := driver.NewCiscoNXOSDevice(d)
			allDevices = append(allDevices, nd)
		case "cisco_smb":
			nd := driver.NewCiscoSMBDevice(d)
			allDevices = append(allDevices, nd)
		case "juniper_junos":
			nd := driver.NewJuniperJunosDevice(d)
			allDevices = append(allDevices, nd)
		default:
			logger.Warningf("device: %s with vendor: %s and platform: %s not supported", d.Name, d.Vendor, d.Platform)
		}
	}

	if !cliParams.NoOp {

		results := []data.Result{}

		var wg sync.WaitGroup
		ch := make(chan data.Result)
		defer close(ch)

		wg.Add(len(allDevices))
		for _, dev := range allDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go driver.RunWithSSH(dev, cliParams.Commands.Commands, ch, &wg)
			case "telnet":
				go driver.RunWithTelnet(dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		for i := 0; i < len(allDevices); i++ {
			results = append(results, <-ch)
		}

		wg.Wait()

		t, err := template.New("results").Parse(templates.CliResult)
		if err != nil {
			logger.Fatal(err)
		}

		fmt.Print(terminal.Banner("Job Results"))

		for _, r := range results {
			if err := t.Execute(os.Stdout, r); err != nil {
				logger.Fatal(err)
			}
		}

		core.WriteToFile(results)
		core.WriteToJSONFile(results)
	}

}
