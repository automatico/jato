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
	commands := []struct {
		command string
		expect  string
	}{
		{
			command: "terminal length 0",
			expect:  "#",
		}, {
			command: "show version",
			expect:  "#",
		}, {
			command: "show ip interface brief",
			expect:  "#",
		}, {
			command: "show cdp neighbors",
			expect:  "#",
		}, {
			command: "show ip arp",
			expect:  "#",
		}, {
			command: "show running-config",
			expect:  "#",
		},
	}

	conn, err := net.Dial("tcp", "192.168.255.150:23")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	auth(conn)

	for _, cmd := range commands {
		// fmt.Println(cmd)
		// fmt.Fprintf(conn, cmd+"\n")
		result := bufferReader(conn, cmd.command, cmd.expect)
		fmt.Println("-------------------------")
		fmt.Println(result)
		fmt.Println("-------------------------")
	}
}

func auth(conn net.Conn) {
	commands := map[string]string{
		"":        "Username:",
		"admin":   "Password:",
		"Juniper": "#",
	}
	var result []string
	for cmd, expect := range commands {
		// fmt.Println(cmd)
		// fmt.Fprintf(conn, cmd+"\n")
		result = append(result, bufferReader(conn, cmd, expect))
	}
	fmt.Println("-------------------------")
	fmt.Println(strings.Join(result, ""))
	fmt.Println("-------------------------")
}

func bufferReader(conn net.Conn, cmd string, expect string) string {
	maxQueueLength := len(expect)
	t := 5 * time.Second
	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 1)       // using small tmo buffer for demonstrating
	queue := []string{}

	fmt.Fprintf(conn, cmd+"\n")

	for {
		conn.SetReadDeadline(time.Now().Add(t))
		n, err := conn.Read(tmp)

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		buf = append(buf, tmp[:n]...)
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
	return string(buf)
}
