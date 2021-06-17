package utils

import (
	"html/template"
	"strings"

	"github.com/automatico/jato/internal/logger"
)

// Underscorer converts a string to an underscore string
// replacing spaces and dashes with underscores
func Underscorer(s string) string {
	re := strings.NewReplacer(" ", "_", "-", "_")
	return re.Replace(s)
}

// CleanOutput removes the first and last lines from
// a string. Strings are split on '\r\n' line endings
func CleanOutput(s string) string {
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
		logger.Fatal("error loading template: %s", err)
	}
	return t
}
