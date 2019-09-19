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
