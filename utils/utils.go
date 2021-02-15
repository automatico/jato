package utils

import "strings"

// Converts a string to an underscore string
// replacing spaces and dashes with underscores
func Underscorer(s string) string {
	re := strings.NewReplacer(" ", "_", "-", "_")
	return re.Replace(s)
}
