package jato

import "regexp"

var (
	UsernameRE *regexp.Regexp = regexp.MustCompile(`(?im)^username:$`)
	PasswordRE *regexp.Regexp = regexp.MustCompile(`(?im)^password:$`)
)
