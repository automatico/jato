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

var aristaEOSDevices []driver.AristaEOSDevice
var arubaAOSCXDevices []driver.ArubaAOSCXDevice
var ciscoAireOSDevices []driver.CiscoAireOSDevice
var ciscoIOSDevices []driver.CiscoIOSDevice
var ciscoIOSXRDevices []driver.CiscoIOSXRDevice
var ciscoNXOSDevices []driver.CiscoNXOSDevice
var ciscoSMBDevices []driver.CiscoSMBDevice
var juniperJunosDevices []driver.JuniperJunosDevice

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
		case "arista_eos":
			ad := driver.NewAristaEOSDevice(d)
			aristaEOSDevices = append(aristaEOSDevices, ad)
		case "aruba_aoscx":
			ad := driver.NewArubaAOSCXDevice(d)
			arubaAOSCXDevices = append(arubaAOSCXDevices, ad)
		case "cisco_aireos":
			cd := driver.NewCiscoAireOSDevice(d)
			creds := data.GetCredentials(cd.Variables.Credentials)
			cd.Credentials = creds
			ciscoAireOSDevices = append(ciscoAireOSDevices, cd)
		case "cisco_ios":
			cd := driver.NewCiscoIOSDevice(d)
			ciscoIOSDevices = append(ciscoIOSDevices, cd)
		case "cisco_iosxr":
			cd := driver.NewCiscoIOSXRDevice(d)
			ciscoIOSXRDevices = append(ciscoIOSXRDevices, cd)
		case "cisco_nxos":
			cd := driver.NewCiscoNXOSDevice(d)
			ciscoNXOSDevices = append(ciscoNXOSDevices, cd)
		case "cisco_smb":
			cd := driver.NewCiscoSMBDevice(d)
			ciscoSMBDevices = append(ciscoSMBDevices, cd)
		case "juniper_junos":
			jd := driver.NewJuniperJunosDevice(d)
			juniperJunosDevices = append(juniperJunosDevices, jd)
		default:
			logger.Warning(fmt.Sprintf("device: %s with vendor: %s and platform: %s not supported", d.Name, d.Vendor, d.Platform))
		}
	}

	if !cliParams.NoOp {

		results := []data.Result{}

		var wg sync.WaitGroup
		ch := make(chan data.Result)
		defer close(ch)

		wg.Add(len(aristaEOSDevices))
		for _, dev := range aristaEOSDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		wg.Add(len(arubaAOSCXDevices))
		for _, dev := range arubaAOSCXDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		wg.Add(len(ciscoAireOSDevices))
		for _, dev := range ciscoAireOSDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

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

		wg.Add(len(ciscoIOSXRDevices))
		for _, dev := range ciscoIOSXRDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		wg.Add(len(ciscoNXOSDevices))
		for _, dev := range ciscoNXOSDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		wg.Add(len(ciscoSMBDevices))
		for _, dev := range ciscoSMBDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		wg.Add(len(juniperJunosDevices))
		for _, dev := range juniperJunosDevices {
			dev := dev // lock the host or the same host can run more than once
			switch dev.Connector {
			case "ssh":
				go network.RunWithSSH(&dev, cliParams.Commands.Commands, ch, &wg)
			}
		}

		devTotal := len(aristaEOSDevices) +
			len(arubaAOSCXDevices) +
			len(ciscoAireOSDevices) +
			len(ciscoIOSDevices) +
			len(ciscoIOSXRDevices) +
			len(ciscoNXOSDevices) +
			len(ciscoSMBDevices) +
			len(juniperJunosDevices)
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
