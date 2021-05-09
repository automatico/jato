package jato

import "os"

// UserCredentials represents a users credentials
type Credentials struct {
	Username      string
	Password      string
	SSHKeyFile    string
	SuperPassword string
}

// Load method to populate a users credentials from
// environment variables.
func (uc Credentials) Load() Credentials {
	uc.Username = os.Getenv("JATO_USERNAME")
	uc.Password = os.Getenv("JATO_PASSWORD")
	uc.SSHKeyFile = os.Getenv("JATO_SSH_KEY_FILE")
	uc.SuperPassword = os.Getenv("JATO_SUPER_PASSWORD")
	return uc
}
