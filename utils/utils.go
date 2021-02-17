package utils

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
