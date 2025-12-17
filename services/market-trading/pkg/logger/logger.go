package logger

import (
	"fmt"
	"log"
	"time"
)

// Info logs an info message
func Info(message string) {
	log.Printf("[INFO] %s | %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// Error logs an error message
func Error(message string) {
	log.Printf("[ERROR] %s | %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// Warn logs a warning message
func Warn(message string) {
	log.Printf("[WARN] %s | %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// Debug logs a debug message
func Debug(message string) {
	log.Printf("[DEBUG] %s | %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// Fatal logs a fatal error and exits
func Fatal(message string) {
	log.Fatalf("[FATAL] %s | %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...))
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	Debug(fmt.Sprintf(format, args...))
}

// Fatalf logs a formatted fatal error and exits
func Fatalf(format string, args ...interface{}) {
	Fatal(fmt.Sprintf(format, args...))
}
