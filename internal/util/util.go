package util

import (
	"html/template"
	"strings"

	"github.com/automatico/jato/internal/logger"
)

// Underscorer converts a string to an underscore string
// replacing spaces and dashes with underscores
func Underscorer(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	// s = strings.ReplaceAll(s, "|", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return s
}

// TruncateOutput removes the first and last lines from
// a string. Strings are split on '\r\n' line endings
func TruncateOutput(s string) string {
	slice := strings.Split(s, "\r\n")
	if len(slice) <= 1 {
		return s
	}
	middle := slice[1 : len(slice)-1]
	return strings.Join(middle, "\r\n")
}

// LoadTemplate loads a template from a filename string
func LoadTemplate(s string) *template.Template {
	t, err := template.ParseFiles(s)
	if err != nil {
		logger.Fatalf("error loading template: %s, %s", s, err)
	}
	return t
}
