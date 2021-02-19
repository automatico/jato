package telnet

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type CommandExpect struct {
	command string
	expect  string
}

// https://stackoverflow.com/questions/24339660/read-whole-data-with-golang-net-conn-read
// https://stackoverflow.com/questions/23531891/how-do-i-succinctly-remove-the-first-element-from-a-slice-in-go
func Telnet() {
	commands := []CommandExpect{
		{"terminal length 0", "#"},
		{"show version", "#"},
		{"show ip interface brief", "#"},
		{"show cdp neighbors", "#"},
		{"show ip arp", "#"},
		{"show running-config", "#"},
	}

	conn, err := net.Dial("tcp", "192.168.255.150:23")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	auth(conn)

	for _, cmd := range commands {
		// fmt.Println("Command:", cmd.command, "Expect:", cmd.expect)
		// fmt.Fprintf(conn, cmd+"\n")
		result := bufferReader(conn, cmd.command, cmd.expect)
		fmt.Println("-------------------------")
		fmt.Println(result)
		fmt.Println("-------------------------")
	}
}

func auth(conn net.Conn) {
	commands := []CommandExpect{
		{"", "Username:"},
		{"admin", "Password:"},
		{"Juniper", "#"},
	}
	var result []string
	for _, cmd := range commands {
		// fmt.Println(cmd)
		// fmt.Fprintf(conn, cmd+"\n")
		result = append(result, bufferReader(conn, cmd.command, cmd.expect))
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
