package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLogger     *log.Logger
	warningLogger  *log.Logger
	warningfLogger *log.Logger
	errorLogger    *log.Logger
	errorfLogger   *log.Logger
	fatalLogger    *log.Logger
	fatalfLogger   *log.Logger
)

func init() {
	infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime)
	warningfLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	errorfLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	fatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime)
	fatalfLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime)
}

func Debug(s string) {
	fmt.Printf("DEBUG => %s\n", s)
}

func Info(v ...interface{}) {
	infoLogger.Println(fmt.Sprint(v...))
}

func Warning(v ...interface{}) {
	warningLogger.Println(fmt.Sprint(v...))
}

func Warningf(format string, v ...interface{}) {
	warningfLogger.Println(fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	errorLogger.Println(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	errorfLogger.Println(fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	fatalLogger.Println(fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	fatalfLogger.Println(fmt.Sprintf(format, v...))
	os.Exit(1)
}
