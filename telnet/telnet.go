package telnet

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// https://stackoverflow.com/questions/24339660/read-whole-data-with-golang-net-conn-read
func Telnet() {
	commands := []string{
		// "\n",
		"admin",
		"Juniper",
		// "terminal length 0",
		// "show version",
		// "show ip interface brief",
		// "show cdp neighbors",
		// "show ip arp",
		// "exit",
	}

	conn, err := net.Dial("tcp", "192.168.255.150:23")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	for _, cmd := range commands {
		// fmt.Println(cmd)
		fmt.Fprintf(conn, cmd+"\n")
	}

	// Works
	// result, _ := ioutil.ReadAll(conn)
	// fmt.Println(string(result))

	// Works
	t := 1 * time.Second
	conn.SetReadDeadline(time.Now().Add(t))
	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 1)       // using small tmo buffer for demonstrating
	container := []string{}
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		buf = append(buf, tmp[:n]...)
		// fmt.Println(string(tmp))
		maxLength := len("Password:")
		if len(container) == maxLength {
			fmt.Println("maxed out")
			fmt.Println(strings.Join(container, ""))
			if strings.Join(container, "") == "Password:" {
				fmt.Println(container)
				fmt.Println(strings.Join(container, ""))
				break
			} else {
				// Pop the front elememnt and shift the rest of the
				// elements left.
				_, container = container[0], container[1:]
				container = append(container, string(tmp))
				fmt.Println("shifted")
				fmt.Println(strings.Join(container, ""))
			}
		} else {
			container = append(container, string(tmp))
			fmt.Println(container)
		}

		switch {
		case strings.Contains(string(tmp[:n]), "Username:"):
			fmt.Println("OHKAY - Username")
			break
			//fmt.Fprintf(conn, "admin"+"\n")
		case strings.Contains(string(tmp[:n]), "Password:"):
			// fmt.Fprintf(conn, "Juniper"+"\n")
			fmt.Println("OHKAY - Password")
			break
		case strings.HasSuffix(string(tmp[:n]), "#"):
			// fmt.Fprintf(conn, "terminal length 0"+"\n")
			fmt.Println("OHKAY - Prompt")
			break

		default:
			fmt.Println("NOOPE")
		}

	}
	fmt.Println("total size:", len(buf))
	fmt.Println(string(buf))

}

func reader(c net.Conn, cmd string) {
	fmt.Fprintf(c, cmd+"\n")

	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
	for {
		n, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		buf = append(buf, tmp[:n]...)
		fmt.Println(string(buf))

	}
	fmt.Println("total size:", len(buf))
	fmt.Println(string(buf))
}

func bufferReader(c net.Conn, expect string, cmd string) {
	fmt.Fprintf(c, cmd+"\n")
	//fmt.Fprintf(conn, "show version")

	buf := make([]byte, 0, 1024) // big buffer
	tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
	for {
		n, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		// fmt.Println("got", n, "bytes.")
		buf = append(buf, tmp[:n]...)
		// fmt.Println(string(buf))
		if strings.HasSuffix(string(buf), expect) {
			fmt.Println(string(buf))
			break
		}

	}
	fmt.Println("total size:", len(buf))
	fmt.Println(string(buf))
}
