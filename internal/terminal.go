package internal

import (
	"fmt"
	"strings"
)

// var termWidth = seperator("#")
const termWidth = "!----------------------------!"

func spacer(n int) string {
	var str strings.Builder
	for i := 0; i < n; i++ {
		str.WriteString(" ")
	}
	return str.String()
}

// Divider is used to output a dividing string
// between outputs. EG:
// ##################
// Job Parameters
// ##################
func Divider(message string) string {
	mlen := len(message)
	spaces := spacer((29 - mlen) / 2)
	return fmt.Sprintf("%s\n!%s%s\n%s\n", termWidth, spaces, message, termWidth)
}
