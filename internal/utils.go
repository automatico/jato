package internal

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
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
	middle := slice[1 : len(slice)-1]
	return strings.Join(middle, "\r\n")
}

// LoadTemplate loads a template from a filename string
func LoadTemplate(s string) *template.Template {
	t, err := template.ParseFiles(s)
	if err != nil {
		panic(err)
	}
	return t
}

// FileStat checks if a file exists and is readable
func FileStat(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Filename: '%s' does not exist or is not readable.", filename)
		os.Exit(1)
	}
}

func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		os.Exit(1)
	}
	return data
}
