package output

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var termWidth = seperator("+")

// CliRunner is the outut for a
// job run from the CLI.
const CliRunner = `{{/* SPACE */}}
{{.divider}}
Credentials:
  - Username:       {{.params.Credentials.Username}}
  - Password:       {{.params.Credentials.Password}}
  - SSH Key File:   {{.params.Credentials.SSHKeyFile}}
  - Super Password: {{.params.Credentials.SuperPassword}}

Devices:
{{- range .params.Devices.Devices }}
  - Name:      {{.Name}}
    IP:        {{.IP}}
    Vendor:    {{.Vendor}}
    Platform:  {{.Platform}}
    Connector: {{.Connector}}
{{- end }}

Commands:
{{- range .params.Commands.Commands}}
  - {{.}}
{{- end }}
{{/* SPACE */}}`

// CliResult is used to display
// the result of a job run
const CliResult = `{{/* SPACE */}}
{{.Device}}:
  OK: {{.OK}}
  Error: {{.Error}}
  Timestamp: {{.Timestamp}}
{{/* SPACE */}}`

// The below is used to determine the size of
// the terminal session.
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWidth() uint {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Col)
}

func seperator(s string) string {
	tw := int(getWidth())
	var str strings.Builder
	for i := 0; i < tw; i++ {
		str.WriteString(s)
	}
	return str.String()
}

// Divider is used to output a dividing string
// between outputs. EG:
// ++++++++++++++++++++
// Job Parameters
// ++++++++++++++++++++
func Divider(message string) string {
	return fmt.Sprintf("%s\n%s\n%s\n", termWidth, message, termWidth)
}
