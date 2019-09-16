package logbuch

import (
	"os"
)

var (
	DebugLogger   = NewLogger(LevelDebug, os.Stdout)
	InfoLogger    = NewLogger(LevelInfo, os.Stdout)
	WarningLogger = NewLogger(LevelWarning, os.Stdout)
	ErrorLogger   = NewLogger(LevelError, os.Stderr)
	level         = LevelDebug
)

func SetLevel(level int) {
	level = getValidLevel(level)
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
	if level <= LevelError {
		ErrorLogger.Log(msg, params...)
	}
}
