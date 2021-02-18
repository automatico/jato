package telnet

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// https://stackoverflow.com/questions/24339660/read-whole-data-with-golang-net-conn-read
func Telnet() {
	commands := []string{
		"\n",
		"terminal length 0",
		"show version",
		"show cdp neighbors",
	}

	conn, err := net.Dial("tcp", "192.168.255.150:23")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	bufferReader(conn, "Username:", "admin")
	bufferReader(conn, "Password:", "Juniper")

	for _, cmd := range commands {
		bufferReader(conn, ">", cmd)
	}
}

func login(c net.Conn) {
	bufferReader(c, "Username:", "\n")
	bufferReader(c, "Password:", "admin")
	bufferReader(c, "#:", "Juniper")
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
