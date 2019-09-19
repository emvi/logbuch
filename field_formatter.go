package logbuch

import (
	"fmt"
	"time"
)

type Fields map[string]interface{}

// FieldFormatter adds fields to the output as key value pairs. The message won't be formatted.
// It prints log messages starting with the timestamp, followed by the log level, the message and key value pairs.
// To make this work the first and only parameter must be of type Fields.
//
// Example:
//  logbuch.Debug("Hello World!", logbuch.Fields{"integer": 123, "string": "test"})
//
// If there is more than one parameter or the type of the parameter is different,
// all parameters will be appended after the message.
type FieldFormatter struct {
	timeFormat string
	separator  string
}

// NewFieldFormatter creates a new FieldFormatter with given timestamp format and separator between message and key value pairs.
func NewFieldFormatter(timeFormat, separator string) *FieldFormatter {
	return &FieldFormatter{timeFormat: timeFormat, separator: separator}
}

// Fmt formats the message as described for the FieldFormatter.
func (formatter *FieldFormatter) Fmt(buffer *[]byte, level int, t time.Time, msg string, params ...interface{}) {
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

	*buffer = append(*buffer, msg...)

	if len(params) > 0 {
		fields, ok := params[0].(Fields)

		if len(params) == 1 && ok {
			*buffer = append(*buffer, formatter.separator...)

			for k, v := range fields {
				*buffer = append(*buffer, fmt.Sprintf(" %s=%v", k, v)...)
			}
		} else {
			*buffer = append(*buffer, formatter.separator...)

			for _, v := range params {
				*buffer = append(*buffer, fmt.Sprintf(" %v", v)...)
			}
		}
	}

	*buffer = append(*buffer, '\n')
}
