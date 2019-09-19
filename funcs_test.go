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

func TestFatal(t *testing.T) {
	var stderr bytes.Buffer

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Fatal must panic")
		}

		if !strings.Contains(stderr.String(), "Fatal message") {
			t.Fatalf("Log must contain error message")
		}
	}()

	SetOutput(nil, &stderr)
	Fatal("Fatal %v", "message")
}

func TestSetFormatter(t *testing.T) {
	formatter := NewFieldFormatter(StandardTimeFormat, "\t")
	SetFormatter(formatter)

	if logger.GetFormatter() != formatter {
		t.Fatal("Formatter must have been set")
	}
}
