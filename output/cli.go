package output

import (
	"strings"
	"syscall"
	"unsafe"
)

var termWidth = seperator()

// JobResult ...
var JobResult = termWidth + "Job Result\n" + termWidth

// CliRunner is the outut for a
// job run from the CLI.
const CliRunner = `{{/* SPACE */}}
--------------------------
Job Parameters
--------------------------
Username: {{.User.Username}}
Password: *************

Devices:
{{- range .Devices.Devices }}
  - Name:      {{.Name}}
    IP:        {{.IP}}
    Vendor:    {{.Vendor}}
    Platform:  {{.Platform}}
    Connector: {{.Connector}}
{{- end }}

Commands:
{{- range .Commands.Commands}}
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

func seperator() string {
	tw := int(getWidth())
	var s strings.Builder
	for i := 1; i < tw; i++ {
		s.WriteString("+")
	}
	return s.String() + "\n"
}
