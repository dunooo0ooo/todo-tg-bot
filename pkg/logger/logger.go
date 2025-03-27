package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func Init(logLevel string) error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	switch logLevel {
	case "debug":
		DebugLogger.SetOutput(file)
	case "info":
		DebugLogger.SetOutput(os.Stderr)
	case "error":
		DebugLogger.SetOutput(os.Stderr)
		InfoLogger.SetOutput(os.Stderr)
	default:
		DebugLogger.SetOutput(os.Stderr)
		InfoLogger.SetOutput(os.Stderr)
	}

	return nil
}

func Info(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	DebugLogger.Printf(format, v...)
}
