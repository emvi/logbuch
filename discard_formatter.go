package logbuch

import (
	"time"
)

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
