package expecter

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Expect struct {
	Command   string `json:"command"`
	Expecting string `json:"expecting"`
	Timeout   int64  `json:"timeout"`
}

type CommandExpect struct {
	CommandExpect []Expect `json:"command_expect"`
}

// Expecter takes a connect, command and string.
// It runs the command against the connection and
// returns the result
func Expecter(connection net.Conn, command string, expecting string, timeout int64) (string, error) {
	// How long to wait for response from device
	// before we giveup and consider it timed out.
	countdown := time.Duration(timeout) * time.Millisecond
	// big buffer holds the result
	buffer := make([]byte, 0, 4096)
	// used to read characters into queue
	tmp := make([]byte, 1)
	// holds number of characters equal to maxQueueLength for
	// matching the expect string
	queue := []string{}
	maxQueueLength := len(expecting)

	// Send command to device
	fmt.Fprintf(connection, command+"\n")

	for {
		// Set timeout for reading from device
		connection.SetReadDeadline(time.Now().Add(countdown))

		n, err := connection.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		buffer = append(buffer, tmp[:n]...)

		if maxQueueLength == 1 && string(tmp) == expecting {
			break
		} else if len(queue) == maxQueueLength {
			if strings.Join(queue, "") == expecting {
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
	return string(buffer), nil
}
