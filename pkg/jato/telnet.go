package jato

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/automatico/jato/internal"
	"github.com/reiver/go-telnet"
)

const TelnetPort int = 23

type TelnetParams struct {
	Port int
}

type TelnetDevice interface {
	ConnectWithTelnet() error
	SendCommandsWithTelnet([]string) Result
	DisconnectTelnet() error
}

func SendCommandsWithTelnet(conn *telnet.Conn, commands []string, expect *regexp.Regexp, timeout int64) ([]CommandOutput, error) {

	cmdOut := []CommandOutput{}

	for _, cmd := range commands {
		res, err := SendCommandWithTelnet(conn, cmd, expect, timeout)
		if err != nil {
			return cmdOut, err
		}
		cmdOut = append(cmdOut, res)
	}

	return cmdOut, nil

}

func SendCommandWithTelnet(conn *telnet.Conn, cmd string, expect *regexp.Regexp, timeout int64) (CommandOutput, error) {

	cmdOut := CommandOutput{}

	writeTelnet(conn, cmd)
	time.Sleep(time.Millisecond * 3)

	res, err := readTelnet(conn, expect, timeout)
	if err != nil {
		return cmdOut, err
	}

	cmdOut.Command = cmd
	cmdOut.CommandU = internal.Underscorer(cmd)
	cmdOut.Output = internal.CleanOutput(res)

	return cmdOut, nil
}

func writeTelnet(w io.Writer, s string) error {
	_, err := w.Write([]byte(s + "\n"))
	if err != nil {
		return err
	}
	return nil
}

func readTelnet(r io.Reader, expect *regexp.Regexp, timeout int64) (string, error) {
	maxBuf := 8192
	tmp := make([]byte, 1)
	result := make([]byte, 0, maxBuf)

	start := time.Now()
	for i := 0; i < maxBuf; i++ {
		if time.Since(start) > time.Second*time.Duration(timeout) {
			err := errors.New("timeout reading from buffer")
			return string(result), err

		}
		n, err := r.Read(tmp)
		if err != nil {
			if err == io.EOF {
				break
			}
			return string(result), err
		}
		result = append(result, tmp[:n]...)
		if expect.MatchString(string(result)) {
			break
		}
	}
	return string(result), nil
}

func RunWithTelnet(td TelnetDevice, commands []string, ch chan Result, wg *sync.WaitGroup) {
	err := td.ConnectWithTelnet()
	if err != nil {
		fmt.Println(err)
	}
	defer td.DisconnectTelnet()
	defer wg.Done()

	result := td.SendCommandsWithTelnet(commands)

	ch <- result
}
