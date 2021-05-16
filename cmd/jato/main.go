package main

import (
	"fmt"
	"html/template"
	"os"
	"sync"
	"time"

	"github.com/automatico/jato/internal"
	"github.com/automatico/jato/internal/templates"
	"github.com/automatico/jato/pkg/jato"
)

var telnetDevices []jato.NetDevice
var sshDevices []jato.NetDevice

func main() {

	timeNow := time.Now().Unix()

	cliParams := jato.CLI()

	jt := jato.Jato{
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

		results := []jato.Result{}

		var wg sync.WaitGroup
		ch := make(chan jato.Result)
		defer close(ch)

		// ssh.SSH(cliParams)
		wg.Add(len(telnetDevices))
		for _, dev := range telnetDevices {
			dev.Credentials.Username = jt.Credentials.Username
			dev.Credentials.Password = jt.Credentials.Password
			go jato.TelnetRunner(dev, jt.CommandExpect, ch, &wg)
		}

		wg.Add(len(sshDevices))
		for _, dev := range sshDevices {
			dev.Credentials.Username = jt.Credentials.Username
			dev.Credentials.Password = jt.Credentials.Password
			dev.SSHParams.Port = 22
			dev.SSHParams.InsecureConnection = true
			dev.SSHParams.InsecureCyphers = true
			fmt.Println(dev)
			go jato.SSHRunner(dev, jt.CommandExpect, ch, &wg)
		}

		devTotal := len(telnetDevices) + len(sshDevices)
		for i := 0; i < devTotal; i++ {
			results = append(results, <-ch)
		}

		wg.Wait()

		t, err := template.New("results").Parse(templates.CliResult)
		if err != nil {
			panic(err)
		}

		fmt.Print(internal.Divider("Job Results"))

		for _, r := range results {
			err = t.Execute(os.Stdout, r)

			if err != nil {
				panic(err)
			}
		}

		jato.WriteToFile(timeNow, results)
		jato.WriteToJSONFile(timeNow, results)
	}

}
