package constant

import "regexp"

var (
	UsernameRE *regexp.Regexp = regexp.MustCompile(`(?im)^username:$`)
	PasswordRE *regexp.Regexp = regexp.MustCompile(`(?im)^password:$`)
)

const (
	SSHPort    int = 22
	TelnetPort int = 23
)
