package jato

import (
	"fmt"
	"net"
	"time"

	"github.com/automatico/jato/command"
	"github.com/automatico/jato/connector"
	"github.com/automatico/jato/device"
	"github.com/automatico/jato/expecter"
	"github.com/automatico/jato/pkg/result"
)

const telnetPort int = 23

// Telnet to a device
func Telnet(jt connector.Jato) result.Results {
	commands := jt.CommandExpect

	results := result.Results{}
	chanResult := make(chan result.Result)
	for _, dev := range jt.Devices.Devices {
		go func(d device.Device, c expecter.CommandExpect) {
			chanResult <- runner(d, c)
		}(dev, commands)
	}

	for range jt.Devices.Devices {
		timeout := time.After(5 * time.Second)
		select {
		case res := <-chanResult:
			results.Results = append(results.Results, res)
			// fmt.Println(res)
		case <-timeout:
			fmt.Println("Timed out!")
		}
	}
	return results
}

func runner(dev device.Device, commands expecter.CommandExpect) result.Result {
	timeNow := time.Now().Unix()
	r := result.Result{}
	r.Device = dev.Name

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dev.IP, telnetPort))
	if err != nil {
		fmt.Println("dial error:", err)
		r.OK = false
		return r
	}
	defer conn.Close()
	auth(conn)

	for _, cmd := range commands.CommandExpect {
		res, err := expecter.Expecter(conn, cmd.Command, cmd.Expecting, cmd.Timeout)
		r.CommandOutputs = append(r.CommandOutputs, result.CommandOutput{Command: cmd.Command, Output: res})
		if err != nil {
			fmt.Println(res)
			// fmt.Println(err)
		}
	}
	r.OK = true
	r.Timestamp = timeNow
	return r
}

func auth(conn net.Conn) {
	commands := []command.CommandExpect{
		{Command: "", Expect: "Username:"},
		{Command: "", Expect: "Username:"},
		{Command: "admin", Expect: "Password:"},
		{Command: "Juniper", Expect: "#"},
	}
	for _, cmd := range commands {
		result, err := expecter.Expecter(conn, cmd.Command, cmd.Expect, 5)
		if err != nil {
			fmt.Println(result)
			fmt.Println(err)
		}
	}
}
