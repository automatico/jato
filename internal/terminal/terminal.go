package terminal

import (
	"fmt"
	"strings"
)

const line = "!----------------------------------------------------------!"

// spacer creates a string of n number of spaces
func spacer(n int) string {
	var str strings.Builder
	for i := 0; i < n; i++ {
		str.WriteString(" ")
	}
	return str.String()
}

// Banner is used to output a banner string
// between outputs. EG:
// !----------------------------------------------------------!
// !                       MESSAGE                            !
// !----------------------------------------------------------!
func Banner(message string) string {
	msgLen := len(message)
	numSpaces := (58 - msgLen) / 2
	preMsgSpaces := spacer(numSpaces)
	var postMsgSpaces string
	if msgLen%2 == 0 {
		postMsgSpaces = spacer(numSpaces)
	} else {
		postMsgSpaces = spacer(numSpaces + 1)
	}
	return fmt.Sprintf("%s\n!%s%s%s!\n%s\n", line, preMsgSpaces, message, postMsgSpaces, line)
}
