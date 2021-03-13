package jato

import (
	"fmt"
	"net"
	"time"
)

const telnetPort int = 23

// Telnet to a device
func Telnet(jt Jato) Results {
	commands := jt.CommandExpect

	results := Results{}
	chanResult := make(chan Result)
	for _, dev := range jt.Devices.Devices {
		go func(d Device, c CommandExpect) {
			chanResult <- runnerTelnet(d, c)
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

func runnerTelnet(dev Device, commands CommandExpect) Result {
	timeNow := time.Now().Unix()
	r := Result{}
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
		res, err := Expecter(conn, cmd.Command, cmd.Expecting, cmd.Timeout)
		r.CommandOutputs = append(r.CommandOutputs, CommandOutput{Command: cmd.Command, Output: res})
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
	commands := CommandExpect{
		[]Expect{
			{Command: "", Expecting: "Username:", Timeout: 5},
			{Command: "", Expecting: "Username:", Timeout: 5},
			{Command: "admin", Expecting: "Password:", Timeout: 5},
			{Command: "Juniper", Expecting: "#", Timeout: 5},
		},
	}
	for _, cmd := range commands.CommandExpect {
		result, err := Expecter(conn, cmd.Command, cmd.Expecting, cmd.Timeout)
		if err != nil {
			fmt.Println(result)
			fmt.Println(err)
		}
	}
}
