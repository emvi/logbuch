package logbuch

import (
	"io"
	"sync"
	"time"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
)

type Logger struct {
	m         sync.Mutex
	level     int
	formatter Formatter
	out       io.Writer
	buffer    []byte

	// PanicOnErr enables panics if the logger cannot write to log output.
	PanicOnErr bool
}

func NewLogger(level int, out io.Writer) *Logger {
	return &Logger{level: getValidLevel(level),
		formatter: NewStandardFormatter(StandardTimeFormat),
		out:       out}
}

func (log *Logger) SetLevel(level int) {
	log.m.Lock()
	defer log.m.Unlock()
	log.level = getValidLevel(level)
}

func (log *Logger) GetLevel() int {
	return log.level
}

func (log *Logger) SetFormatter(formatter Formatter) {
	log.m.Lock()
	defer log.m.Unlock()
	log.formatter = formatter
}

func (log *Logger) GetFormatter() Formatter {
	return log.formatter
}

func (log *Logger) SetOut(out io.Writer) {
	log.m.Lock()
	defer log.m.Unlock()
	log.out = out
}

func (log *Logger) GetOut() io.Writer {
	return log.out
}

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
