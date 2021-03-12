package telnet

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/automatico/jato/command"
	"github.com/automatico/jato/connector"
	"github.com/automatico/jato/device"
	"github.com/automatico/jato/expecter"
	"github.com/automatico/jato/result"
)

const telnetPort int = 23

// Telnet to a device
func Telnet(jt connector.Jato) result.Results {

	//timeNow := time.Now().Unix()
	//usr := cp.User
	//cmds := cp.Commands
	// devs := cp.Devices

	commands := expecter.CommandExpect{
		CommandExpect: []expecter.Expect{
			{Command: "terminal length 0", Expecting: "#", Timeout: 5},
			{Command: "show version", Expecting: "#", Timeout: 5},
			{Command: "show ip interface brief", Expecting: "#", Timeout: 5},
			{Command: "show cdp neighbors", Expecting: "#", Timeout: 5},
			{Command: "show ip arp", Expecting: "#", Timeout: 5},
			{Command: "show running-config", Expecting: "#", Timeout: 5},
		},
	}

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
			break
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
		res := bufferReader(conn, cmd.Command, cmd.Expecting)
		r.CommandOutputs = append(r.CommandOutputs, result.CommandOutput{Command: cmd.Command, Output: res})
	}
	r.OK = true
	r.Timestamp = timeNow
	return r
}

func auth(conn net.Conn) {
	commands := []command.CommandExpect{
		{Command: "", Expect: "Username:"},
		{Command: "admin", Expect: "Password:"},
		{Command: "Juniper", Expect: "#"},
	}
	var res []string
	for _, cmd := range commands {
		res = append(res, bufferReader(conn, cmd.Command, cmd.Expect))
	}
}

func bufferReader(conn net.Conn, cmd string, expect string) string {
	// How long to wait for response from device
	// before we giveup and consider it timed out.
	timeout := 5 * time.Second
	// big buffer holds the result
	buffer := make([]byte, 0, 4096)
	// used to read characters into queue
	tmp := make([]byte, 1)
	// holds number of characters equal to maxQueueLength for
	// matching the expect string
	queue := []string{}
	maxQueueLength := len(expect)

	// Send command to device
	fmt.Fprintf(conn, cmd+"\n")

	for {
		// Set timeout for reading from device
		conn.SetReadDeadline(time.Now().Add(timeout))

		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		buffer = append(buffer, tmp[:n]...)

		if maxQueueLength == 1 && string(tmp) == expect {
			break
		} else if len(queue) == maxQueueLength {
			if strings.Join(queue, "") == expect {
				break
			} else {
				// Pop the front elememnt and shift the rest of the
				// elements left.
				_, queue = queue[0], queue[1:]
				// Add element to the end of the queue
				queue = append(queue, string(tmp))
			}
		} else {
			// Queue is not full, so add elements to queue.
			queue = append(queue, string(tmp))
		}

	}
	return string(buffer)
}
