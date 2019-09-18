package logbuch

import (
	"io"
	"os"
)

var (
	logger = NewLogger(os.Stdout, os.Stderr)
)

// SetOutput sets the output channels for the default logger.
// The first parameter is used for debug, info and warnings.
// The second one for error logs.
func SetOutput(stdout, stderr io.Writer) {
	logger.SetOut(LevelDebug, stdout)
	logger.SetOut(LevelInfo, stdout)
	logger.SetOut(LevelWarning, stdout)
	logger.SetOut(LevelError, stderr)
}

// SetLevel sets the logging level.
func SetLevel(level int) {
	logger.SetLevel(level)
}

// Debug logs a formatted debug message.
func Debug(msg string, params ...interface{}) {
	logger.Debug(msg, params...)
}

// Info logs a formatted info message.
func Info(msg string, params ...interface{}) {
	logger.Info(msg, params...)
}

// Warn logs a formatted warning message.
func Warn(msg string, params ...interface{}) {
	logger.Warn(msg, params...)
}

// Error logs a formatted error message.
func Error(msg string, params ...interface{}) {
	// maximum level cannot be disabled
	logger.Error(msg, params...)
}
