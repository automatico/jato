package constant

import "regexp"

const (
	SSHPort    int = 22
	TelnetPort int = 23
)

var (
	UsernameRE *regexp.Regexp = regexp.MustCompile(`(?im)^username:$`)
	PasswordRE *regexp.Regexp = regexp.MustCompile(`(?im)^password:$`)
)

var InsecureSSHCyphers = []string{
	"aes128-ctr",
	"aes192-ctr",
	"aes256-ctr",
	"aes128-cbc",
	"aes192-cbc",
	"aes256-cbc",
	"3des-cbc",
	"des-cbc",
}

var InsecureSSHKeyAlgorithms = []string{
	"diffie-hellman-group-exchange-sha256",
	"diffie-hellman-group-exchange-sha1",
	"diffie-hellman-group1-sha1",
	"diffie-hellman-group14-sha1",
}
