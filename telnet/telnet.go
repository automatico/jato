package telnet

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/automatico/jato/command"
	"github.com/automatico/jato/device"
)

const telnetPort int = 23

// Telnet to a device
func Telnet() {
	devices := []device.Device{
		{Name: "iosv-1", IP: "192.168.255.150", Vendor: "cisco", Platform: "ios", Connector: "telnet"},
		{Name: "iosv-4", IP: "192.168.255.154", Vendor: "cisco", Platform: "ios", Connector: "telnet"},
		{Name: "iosv-5", IP: "192.168.255.155", Vendor: "cisco", Platform: "ios", Connector: "telnet"},
		{Name: "iosv-6", IP: "192.168.255.156", Vendor: "cisco", Platform: "ios", Connector: "telnet"},
		{Name: "iosv-7", IP: "192.168.255.157", Vendor: "cisco", Platform: "ios", Connector: "telnet"},
	}

	commands := []command.CommandExpect{
		{Command: "terminal length 0", Expect: "#"},
		{Command: "show version", Expect: "#"},
		{Command: "show ip interface brief", Expect: "#"},
		{Command: "show cdp neighbors", Expect: "#"},
		{Command: "show ip arp", Expect: "#"},
		{Command: "show running-config", Expect: "#"},
	}

	results := make(chan bool)
	for _, dev := range devices {
		go func(d device.Device, c []command.CommandExpect) {
			results <- runner(d, c)
		}(dev, commands)
	}

	for range devices {
		timeout := time.After(10 * time.Second)
		select {
		case result := <-results:
			fmt.Println(result)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}

}

func runner(dev device.Device, commands []command.CommandExpect) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dev.IP, telnetPort))
	if err != nil {
		fmt.Println("dial error:", err)
		return false
	}
	defer conn.Close()

	auth(conn)

	for _, cmd := range commands {
		// fmt.Println("Command:", cmd.command, "Expect:", cmd.expect)
		// fmt.Fprintf(conn, cmd+"\n")
		result := bufferReader(conn, cmd.Command, cmd.Expect)
		fmt.Println("-------------------------")
		fmt.Println(result)
		fmt.Println("-------------------------")
	}
	return true
}

func auth(conn net.Conn) {
	commands := []command.CommandExpect{
		{Command: "", Expect: "Username:"},
		{Command: "admin", Expect: "Password:"},
		{Command: "Juniper", Expect: "#"},
	}
	var result []string
	for _, cmd := range commands {
		// fmt.Println(cmd)
		// fmt.Fprintf(conn, cmd+"\n")
		result = append(result, bufferReader(conn, cmd.Command, cmd.Expect))
	}
	fmt.Println("-------------------------")
	fmt.Println(strings.Join(result, ""))
	fmt.Println("-------------------------")
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
				queue = append(queue, string(tmp))
			}
		} else {
			// Queue is not full, so add elements to queue.
			queue = append(queue, string(tmp))
		}

	}
	return string(buffer)
}
