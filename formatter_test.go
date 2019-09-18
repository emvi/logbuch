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
		formatter.Fmt(&buffer, in.level, now, in.msg, in.params...)

		if string(buffer) != expected[i] {
			t.Fatalf("Expected '%v' but was: %v", expected[i], string(buffer))
		}
	}
}

func TestDiscardFormatter(t *testing.T) {
	formatter := NewDiscardFormatter()
	now := time.Now()
	var buffer []byte
	formatter.Fmt(&buffer, LevelDebug, now, "Message", "param", 123)

	if string(buffer) != "" {
		t.Fatalf("Log must be discarded but was: %v", string(buffer))
	}
}
