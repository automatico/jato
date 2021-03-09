package credentials

import "os"

type UserCredentials struct {
	Username      string
	Password      string
	SSHKeyFile    string
	SuperPassword string
}

func LoadUser() UserCredentials {
	cred := UserCredentials{
		Username:      os.Getenv("JATO_USERNAME"),
		Password:      os.Getenv("JATO_PASSWORD"),
		SSHKeyFile:    os.Getenv("JATO_SSH_KEY_FILE"),
		SuperPassword: os.Getenv("JATO_SUPER_PASSWORD"),
	}
	return cred
}
