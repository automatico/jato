package logger

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {

	InfoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

func Info(s string) {
	InfoLogger.Println(s)
}

func Warning(s string) {
	WarningLogger.Println(s)
}

func Error(s string) {
	ErrorLogger.Println(s)
}
