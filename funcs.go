package logbuch

import (
	"io"
	"os"
)

var (
	DebugLogger   = NewLogger(LevelDebug, os.Stdout)
	InfoLogger    = NewLogger(LevelInfo, os.Stdout)
	WarningLogger = NewLogger(LevelWarning, os.Stdout)
	ErrorLogger   = NewLogger(LevelError, os.Stderr)
	level         = LevelDebug
)

func SetOutput(stdout, stderr io.Writer) {
	DebugLogger.SetOut(stdout)
	InfoLogger.SetOut(stdout)
	WarningLogger.SetOut(stdout)
	ErrorLogger.SetOut(stderr)
}

func SetLevel(lvl int) {
	level = getValidLevel(lvl)
}

func Debug(msg string, params ...interface{}) {
	if level <= LevelDebug {
		DebugLogger.Log(msg, params...)
	}
}

func Info(msg string, params ...interface{}) {
	if level <= LevelInfo {
		InfoLogger.Log(msg, params...)
	}
}

func Warn(msg string, params ...interface{}) {
	if level <= LevelWarning {
		WarningLogger.Log(msg, params...)
	}
}

func Error(msg string, params ...interface{}) {
	// maximum level cannot be disabled
	ErrorLogger.Log(msg, params...)
}
