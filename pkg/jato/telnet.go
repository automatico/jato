package jato

import (
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/automatico/jato/internal"
)

func TelnetExpecter(conn net.Conn, cmd string, expecting string, timeout int64) (string, error) {
	// How long to wait for response from device
	// before we give up and consider it timed out.
	countdown := time.Duration(timeout) * time.Second

	// Big buffer holds the result
	result := make([]byte, 0, 4096)
	// Used to read characters into queue
	tmp := make([]byte, 1)

	// Holds number of characters equal to maxQueueLength for
	// matching the expect string
	queue := []string{}
	maxQueueLength := len(expecting)

	// Send command to device
	// fmt.Fprintf(conn, cmd+"\n")
	conn.Write([]byte(cmd + "\n"))

	for {
		// Set timeout for reading from device
		conn.SetReadDeadline(time.Now().Add(countdown))

		n, err := conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				break // Reached the end of file
			}
			return "read error", err
		}

		result = append(result, tmp[:n]...)

		if maxQueueLength == 1 && string(tmp) == expecting {
			// Your done, exit the loop
			break
		} else if len(queue) == maxQueueLength {
			// Queue is full, check for expecting string
			if strings.Join(queue, "") == expecting {
				// Your done, exit the loop
				break
			} else {
				// Pop the front elememnt and shift the rest of the
				// elements left
				_, queue = queue[0], queue[1:]
				// Add element to the end of the queue
				queue = append(queue, string(tmp))
			}
		} else {
			// Queue is not full, so add elements to queue
			queue = append(queue, string(tmp))
		}

	}
	return string(result), nil
}

func TelnetRunner(nd NetDevice, ce CommandExpect, ch chan Result, wg *sync.WaitGroup) {
	conn := nd.ConnectWithTelnet()
	defer conn.Close()
	defer wg.Done()

	result := Result{}
	cmdOut := []CommandOutput{}

	result.Device = nd.Name
	result.Timestamp = time.Now().Unix()
	for _, cmd := range ce.CommandExpect {
		res, err := TelnetExpecter(conn, cmd.Command, cmd.Expecting, cmd.Timeout)
		if err != nil {
			result.OK = false
			ch <- result
			return
		}
		out := CommandOutput{Command: internal.Underscorer(cmd.Command), Output: res}
		cmdOut = append(cmdOut, out)
	}
	result.CommandOutputs = cmdOut
	result.OK = true
	ch <- result
}
