package constant

import (
	"os"
	"path/filepath"
	"regexp"
)

const SSHPort = 22
const TelnetPort = 23

const Timeout = 5

var LoginRE = regexp.MustCompile(`(?im)^login:$`)
var UsernameRE = regexp.MustCompile(`(?im)^username:$`)
var PasswordRE = regexp.MustCompile(`(?im)^password:$`)

var SSHKnownHostsFile = filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
var SSHKeyFile = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")

var InsecureSSHCiphers = []string{
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
