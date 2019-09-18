package logbuch

import (
	"fmt"
	"time"
)

const (
	// StandardTimeFormat is a synonym for time.RFC3339Nano.
	StandardTimeFormat = time.RFC3339Nano
)

// Formatter is an interface to format log messages.
type Formatter interface {
	// Fmt formats a logger message and writes the result into the buffer.
	Fmt(*[]byte, int, time.Time, string, ...interface{})
}

// StandardFormatter is the default logger.
// It prints log messages starting with the timestamp, followed by the log level and the formatted message.
type StandardFormatter struct {
	timeFormat string
}

// NewStandardFormatter creates a new StandardFormatter with given timestamp format.
func NewStandardFormatter(timeFormat string) *StandardFormatter {
	return &StandardFormatter{timeFormat: timeFormat}
}

// Fmt formats the message as described for the StandardFormatter.
func (formatter *StandardFormatter) Fmt(buffer *[]byte, level int, t time.Time, msg string, params ...interface{}) {
	*buffer = append(*buffer, t.Format(formatter.timeFormat)+" "...)

	switch level {
	case LevelDebug:
		*buffer = append(*buffer, "[DEBUG] "...)
	case LevelInfo:
		*buffer = append(*buffer, "[INFO ] "...)
	case LevelWarning:
		*buffer = append(*buffer, "[WARN ] "...)
	case LevelError:
		*buffer = append(*buffer, "[ERROR] "...)
	}

	if len(params) == 0 {
		*buffer = append(*buffer, msg...)
	} else {
		*buffer = append(*buffer, fmt.Sprintf(msg, params...)...)
	}

	if len(*buffer) == 0 || (*buffer)[len(*buffer)-1] != '\n' {
		*buffer = append(*buffer, '\n')
	}
}

// DiscardFormatter drops all messages.
type DiscardFormatter struct{}

// NewDiscardFormatter creates a new DiscardFormatter.
func NewDiscardFormatter() *DiscardFormatter {
	return new(DiscardFormatter)
}

// Fmt drops the message.
func (formatter *DiscardFormatter) Fmt(buffer *[]byte, level int, t time.Time, msg string, params ...interface{}) {
	// does nothing
}
