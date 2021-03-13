package jato

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
)

// Expect struct
// Command: command to run
// Expecting: string you are expecting
// Timeout: How long to wait for a command
type Expect struct {
	Command   string `json:"command"`
	Expecting string `json:"expecting"`
	Timeout   int64  `json:"timeout"`
}

// CommandExpect holds a slice of
// Expect Structs
type CommandExpect struct {
	CommandExpect []Expect `json:"command_expect"`
}

// Expecter takes a connect, command and string.
// It runs the command against the connection and
// returns the result
func Expecter(connection net.Conn, command string, expecting string, timeout int64) (string, error) {
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
	fmt.Fprintf(connection, command+"\n")

	for {
		// Set timeout for reading from device
		connection.SetReadDeadline(time.Now().Add(countdown))

		n, err := connection.Read(tmp)
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

func LoadCommands(fileName string) CommandExpect {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	data := CommandExpect{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
