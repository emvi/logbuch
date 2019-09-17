package logbuch

import (
	"io"
	"os"
)

var (
	// DebugLogger is used for debug logs.
	// Logs to os.Stdout by default.
	DebugLogger = NewLogger(LevelDebug, os.Stdout)

	// InfoLogger is used for info logs.
	// Logs to os.Stdout by default.
	InfoLogger = NewLogger(LevelInfo, os.Stdout)

	// WarningLogger is used for warning logs.
	// Logs to os.Stdout by default.
	WarningLogger = NewLogger(LevelWarning, os.Stdout)

	// ErrorLogger is used for error logs.
	// Logs to os.Stderr by default.
	ErrorLogger = NewLogger(LevelError, os.Stderr)

	level = LevelDebug
)

// SetOutput sets the output channels for the default loggers.
// The first parameter is used for debug, info and warnings.
// The second one for error logs.
func SetOutput(stdout, stderr io.Writer) {
	DebugLogger.SetOut(stdout)
	InfoLogger.SetOut(stdout)
	WarningLogger.SetOut(stdout)
	ErrorLogger.SetOut(stderr)
}

// SetLevel sets the logging level.
func SetLevel(lvl int) {
	level = getValidLevel(lvl)
}

// Debug logs a formatted debug message.
func Debug(msg string, params ...interface{}) {
	if level <= LevelDebug {
		DebugLogger.Log(msg, params...)
	}
}

// Info logs a formatted info message.
func Info(msg string, params ...interface{}) {
	if level <= LevelInfo {
		InfoLogger.Log(msg, params...)
	}
}

// Warn logs a formatted warning message.
func Warn(msg string, params ...interface{}) {
	if level <= LevelWarning {
		WarningLogger.Log(msg, params...)
	}
}

// Error logs a formatted error message.
func Error(msg string, params ...interface{}) {
	// maximum level cannot be disabled
	ErrorLogger.Log(msg, params...)
}
