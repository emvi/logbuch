package logbuch

import (
	"testing"
	"time"
)

func TestDiscardFormatter(t *testing.T) {
	formatter := NewDiscardFormatter()
	now := time.Now()
	var buffer []byte
	formatter.Fmt(&buffer, LevelDebug, now, "Message", "param", 123)

	if string(buffer) != "" {
		t.Fatalf("Log must be discarded but was: %v", string(buffer))
	}
}
