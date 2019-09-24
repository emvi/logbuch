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

func TestDiscardFormatterPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Formatter must panic")
		} else {
			if r != "message" {
				t.Fatalf("Message not correct: %v", r)
			}
		}
	}()

	formatter := NewDiscardFormatter()
	formatter.Pnc("message")
}

func TestDiscardFormatterPanicFmt(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Formatter must panic")
		} else {
			if r != "message formatted" {
				t.Fatalf("Message must be formatted, but was: %v", r)
			}
		}
	}()

	formatter := NewDiscardFormatter()
	formatter.Pnc("message %s", "formatted")
}
