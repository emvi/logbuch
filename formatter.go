package logbuch

import (
	"fmt"
	"time"
)

const (
	StandardTimeFormat = time.RFC3339Nano
)

type Formatter interface {
	Fmt(*[]byte, int, time.Time, string, ...interface{})
}

type StandardFormatter struct {
	timeFormat string
}

func NewStandardFormatter(timeFormat string) *StandardFormatter {
	return &StandardFormatter{timeFormat: timeFormat}
}

func (formatter *StandardFormatter) Fmt(buffer *[]byte, level int, t time.Time, msg string, params ...interface{}) {
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

	*buffer = append(*buffer, t.Format(formatter.timeFormat)+" "...)
	*buffer = append(*buffer, fmt.Sprintf(msg, params...)...)
}
