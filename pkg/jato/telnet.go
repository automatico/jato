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
			chanResult <- telnetRunner(d, c)
		}(dev, commands)
	}

	for range jt.Devices.Devices {
		select {
		case res := <-chanResult:
			results.Results = append(results.Results, res)
			// fmt.Println(res)
		case <-time.After(6 * time.Second):
			fmt.Println("Timed out!")
		}
	}
	return results
}

func telnetRunner(dev Device, commands CommandExpect) Result {
	timeNow := time.Now().Unix()
	r := Result{}
	r.Device = dev.Name
	r.Timestamp = timeNow

	// fmt.Println("DIALING")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dev.IP, telnetPort))
	if err != nil {
		fmt.Println("dial error:", err)
		r.OK = false
		return r
	}
	defer conn.Close()
	// fmt.Println("FINISHED DIALING")

	auth(conn)

	for _, cmd := range commands.CommandExpect {
		res, err := Expecter(conn, cmd.Command, cmd.Expecting, cmd.Timeout)
		r.CommandOutputs = append(r.CommandOutputs, CommandOutput{Command: cmd.Command, Output: res})
		if err != nil {
			fmt.Println(res)
			fmt.Println(err)
		}
	}

	theTimeNow := time.Now().Unix()
	// fmt.Println(theTimeNow)
	// fmt.Println(timeNow + 5)
	if theTimeNow > timeNow+5 {
		// Consider the device timed out sending commands
		// fmt.Println("Timeout waiting for commands")
		r.OK = false
		r.Error = "timeout sending commands"
	} else {
		r.OK = true
	}
	return r
}

func auth(conn net.Conn) {
	commands := CommandExpect{
		[]Expect{
			{Command: "", Expecting: "Username:", Timeout: 2},
			{Command: "admin", Expecting: "Password:", Timeout: 2},
			{Command: "Juniper", Expecting: "#", Timeout: 2},
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