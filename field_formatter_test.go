package logbuch

import (
	"strings"
	"testing"
	"time"
)

func TestFieldFormatter(t *testing.T) {
	formatter := NewFieldFormatter(StandardTimeFormat, "\t\t\t")
	now := time.Now()
	nowStr := now.Format(StandardTimeFormat)
	var buffer []byte
	input := []struct {
		level  int
		msg    string
		params []interface{}
	}{
		{LevelDebug, "Hello World!", nil},
		{LevelInfo, "Hello World!", []interface{}{"test"}},
		{LevelWarning, "Hello World!", []interface{}{Fields{"text": "test", "integer": 123}}},
		{LevelError, "Hello World!", []interface{}{Fields{"text": "test", "integer": 123, "float": -3.14}}},
	}
	expected := [][]string{
		{nowStr + " [DEBUG] " + "Hello World!\n"},
		{nowStr + " [INFO ] " + "Hello World!", "test", "\n"},
		{nowStr + " [WARN ] " + "Hello World!\t\t\t", "text=test", "integer=123", "\n"},
		{nowStr + " [ERROR] " + "Hello World!\t\t\t", "text=test", "integer=123", "float=-3.14", "\n"},
	}

	for i, in := range input {
		buffer = buffer[:0]
		formatter.Fmt(&buffer, in.level, now, in.msg, in.params...)
		out := string(buffer)
		t.Log(out)

		for _, exp := range expected[i] {
			if !strings.Contains(out, exp) {
				t.Fatalf("Expected '%v' in '%v'", exp, out)
			}
		}
	}
}

func TestFieldFormatterPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Formatter must panic")
		} else {
			if r != "message" {
				t.Fatalf("Message not correct: %v", r)
			}
		}
	}()

	formatter := NewFieldFormatter(StandardTimeFormat, "\t\t\t")
	formatter.Pnc("message")
}

func TestFieldFormatterPanicFmt(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Formatter must panic")
		} else {
			if r != "message formatted" {
				t.Fatalf("Message must be formatted, but was: %v", r)
			}
		}
	}()

	formatter := NewFieldFormatter(StandardTimeFormat, "\t\t\t")
	formatter.Pnc("message %s", "formatted")
}

func TestFieldFormatterPanicFmtFields(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Formatter must panic")
		} else {
			rStr, _ := r.(string)

			if !strings.Contains(rStr, "message") ||
				!strings.Contains(rStr, "variable=value") ||
				!strings.Contains(rStr, "more=123") {
				t.Fatalf("Message must be formatted, but was: %v", r)
			}
		}
	}()

	formatter := NewFieldFormatter(StandardTimeFormat, "\t\t\t")
	formatter.Pnc("message", Fields{"variable": "value", "more": 123})
}
