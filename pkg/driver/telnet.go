package driver

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/automatico/jato/internal/utils"
	"github.com/automatico/jato/pkg/constant"
	"github.com/automatico/jato/pkg/data"
	"github.com/reiver/go-telnet"
)

type TelnetParams struct {
	Port int `json:"port"`
}

type TelnetDevice interface {
	ConnectWithTelnet() error
	SendCommandsWithTelnet([]string) data.Result
	DisconnectTelnet() error
}

func InitTelnetParams(s *TelnetParams) {
	if s.Port == 0 {
		s.Port = constant.TelnetPort
	}
}

func SendCommandsWithTelnet(conn *telnet.Conn, commands []string, expect *regexp.Regexp, timeout int64) ([]data.CommandOutput, error) {

	cmdOut := []data.CommandOutput{}

	for _, cmd := range commands {
		res, err := SendCommandWithTelnet(conn, cmd, expect, timeout)
		if err != nil {
			return cmdOut, err
		}
		cmdOut = append(cmdOut, res)
	}

	return cmdOut, nil

}

func SendCommandWithTelnet(conn *telnet.Conn, cmd string, expect *regexp.Regexp, timeout int64) (data.CommandOutput, error) {

	cmdOut := data.CommandOutput{}

	WriteTelnet(conn, cmd)
	time.Sleep(time.Millisecond * 3)

	res, err := ReadTelnet(conn, expect, timeout)
	if err != nil {
		return cmdOut, err
	}

	cmdOut.Command = cmd
	cmdOut.CommandU = utils.Underscorer(cmd)
	cmdOut.Output = utils.CleanOutput(res)

	return cmdOut, nil
}

func WriteTelnet(w io.Writer, s string) error {
	_, err := w.Write([]byte(s + "\n"))
	if err != nil {
		return err
	}
	return nil
}

func ReadTelnet(r io.Reader, expect *regexp.Regexp, timeout int64) (string, error) {
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
		// Uncomment for testing
		// fmt.Print(string(tmp))
		result = append(result, tmp[:n]...)
		if expect.MatchString(string(result)) {
			break
		}
	}
	return string(result), nil
}

func RunWithTelnet(nd NetDevice, commands []string, ch chan data.Result, wg *sync.WaitGroup) {
	err := nd.ConnectWithTelnet()
	if err != nil {
		fmt.Println(err)
	}
	defer nd.DisconnectTelnet()
	defer wg.Done()

	result := nd.SendCommandsWithTelnet(commands)

	ch <- result
}
