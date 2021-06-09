package logger

import (
	"log"
	"os"
)

var (
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
)

func init() {

	infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

func Info(s string) {
	infoLogger.Println(s)
}

func Warning(s string) {
	warningLogger.Println(s)
}

func Error(s string) {
	errorLogger.Println(s)
}
