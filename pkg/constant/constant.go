package constant

import "regexp"

var (
	LoginRE    *regexp.Regexp = regexp.MustCompile(`(?im)login:`)
	UsernameRE *regexp.Regexp = regexp.MustCompile(`(?im)^username:\s$`)
	PasswordRE *regexp.Regexp = regexp.MustCompile(`(?im)^password:\s$`)
)

const (
	SSHPort    int = 22
	TelnetPort int = 23
)
