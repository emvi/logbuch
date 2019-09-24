package logbuch

import (
	"testing"
	"time"
)

func TestStandardFormatter(t *testing.T) {
	formatter := NewStandardFormatter(StandardTimeFormat)
	now := time.Now()
	nowStr := now.Format(StandardTimeFormat)
	var buffer []byte
	input := []struct {
		level  int
		msg    string
		params []interface{}
	}{
		{LevelDebug, "Hello World!", nil},
		{LevelInfo, "Hello %s!", []interface{}{"World"}},
		{LevelWarning, "Hello %s %d!", []interface{}{"World", 123}},
		{LevelError, "Hello %s %d %v!", []interface{}{"World", 123, -3.14}},
	}
	expected := []string{
		nowStr + " [DEBUG] " + "Hello World!\n",
		nowStr + " [INFO ] " + "Hello World!\n",
		nowStr + " [WARN ] " + "Hello World 123!\n",
		nowStr + " [ERROR] " + "Hello World 123 -3.14!\n",
	}

	for i, in := range input {
		buffer = buffer[:0]
		formatter.Fmt(&buffer, in.level, now, in.msg, in.params)
		out := string(buffer)
		t.Log(out)

		if out != expected[i] {
			t.Fatalf("Expected '%v' but was: %v", expected[i], out)
		}
	}
}

func TestStandardFormatterPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Formatter must panic")
		} else {
			if r != "message" {
				t.Fatalf("Message not correct: %v", r)
			}
		}
	}()

	formatter := NewStandardFormatter(StandardTimeFormat)
	formatter.Pnc("message", nil)
}

func TestStandardFormatterPanicFmt(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Formatter must panic")
		} else {
			if r != "message formatted" {
				t.Fatalf("Message must be formatted, but was: %v", r)
			}
		}
	}()

	formatter := NewStandardFormatter(StandardTimeFormat)
	formatter.Pnc("message %s", []interface{}{"formatted"})
}

func TestStandardFormatterDiableTime(t *testing.T) {
	formatter := NewStandardFormatter("")
	var buffer []byte
	formatter.Fmt(&buffer, LevelDebug, time.Now(), "message", nil)

	if string(buffer) != "[DEBUG] message\n" {
		t.Fatalf("Unexpected log: %v", string(buffer))
	}
}
