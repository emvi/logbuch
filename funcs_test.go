package logbuch

import (
	"bytes"
	"strings"
	"testing"
)

func TestFuncs(t *testing.T) {
	for lvl := 0; lvl < LevelError; lvl++ {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		SetOutput(&stdout, &stderr)
		SetLevel(lvl)
		Debug("Debug %s", "message")
		Info("Info %s", "message")
		Warn("Warning %s", "message")
		Error("Error %s", "message")

		if strings.Contains(stdout.String(), "Debug message") != (lvl <= LevelDebug) ||
			strings.Contains(stdout.String(), "Info message") != (lvl <= LevelInfo) ||
			strings.Contains(stdout.String(), "Warning message") != (lvl <= LevelWarning) ||
			strings.Contains(stdout.String(), "Error message") {
			t.Fatalf("Unexpected standard log output for level %v: %v", logger.level, stdout.String())
		}

		if !strings.Contains(stderr.String(), "Error message") {
			t.Fatalf("Unexpected error log output for level %v: %v", logger.level, stderr.String())
		}
	}
}
