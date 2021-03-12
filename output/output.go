package output

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var termWidth = seperator("#")

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
// ##################
// Job Parameters
// ##################
func Divider(message string) string {
	return fmt.Sprintf("%s\n%s\n%s\n", termWidth, message, termWidth)
}
