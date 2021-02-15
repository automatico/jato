package user

import (
	"os"
)

// User represents a users credentials
type User struct {
	Username string
	Password string
}

func LoadUser() User {
	usr := User{
		Username: os.Getenv("JATO_SSH_USER"),
		Password: os.Getenv("JATO_SSH_PASS"),
	}
	return usr
}
