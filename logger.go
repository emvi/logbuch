package logbuch

import (
	"io"
	"sync"
	"time"
)

const (
	// LevelDebug log all messages.
	LevelDebug = iota

	// LevelInfo log info, warning and error messages.
	LevelInfo

	// LevelWarning log warning and error messages.
	LevelWarning

	// LevelError log error messages only.
	LevelError
)

// Logger writes messages to a defined output io.Writer by using a given Formatter.
type Logger struct {
	m         sync.Mutex
	level     int
	formatter Formatter
	out       io.Writer
	buffer    []byte

	// PanicOnErr enables panics if the logger cannot write to log output.
	PanicOnErr bool
}

// NewLogger creates a new logger for given log level and io.Writer that uses the StandardFormatter.
func NewLogger(level int, out io.Writer) *Logger {
	return &Logger{level: getValidLevel(level),
		formatter: NewStandardFormatter(StandardTimeFormat),
		out:       out}
}

// SetLevel sets the log level.
func (log *Logger) SetLevel(level int) {
	log.m.Lock()
	defer log.m.Unlock()
	log.level = getValidLevel(level)
}

// GetLevel returns the log level.
func (log *Logger) GetLevel() int {
	return log.level
}

// SetFormatter sets the formatter.
func (log *Logger) SetFormatter(formatter Formatter) {
	log.m.Lock()
	defer log.m.Unlock()
	log.formatter = formatter
}

// GetFormatter returns the formatter.
func (log *Logger) GetFormatter() Formatter {
	return log.formatter
}

// SetOut sets the io.Writer for log output.
func (log *Logger) SetOut(out io.Writer) {
	log.m.Lock()
	defer log.m.Unlock()
	log.out = out
}

// GetOut returns the io.Writer for log output.
func (log *Logger) GetOut() io.Writer {
	return log.out
}

// Log writes a message with given parameters to log.
// How the messages is formatted depends on the configured formatter.
func (log *Logger) Log(msg string, params ...interface{}) {
	now := time.Now()
	log.m.Lock()
	defer log.m.Unlock()
	log.buffer = log.buffer[:0]
	log.formatter.Fmt(&log.buffer, log.level, now, msg, params...)

	// panic in case the logger cannot write to the configured io.Writer and panic is enabled
	if _, err := log.out.Write(log.buffer); err != nil && log.PanicOnErr {
		panic(err)
	}
}

func getValidLevel(level int) int {
	if level < LevelDebug || level > LevelError {
		return LevelDebug
	}

	return level
}
