package logbuch

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLoggerLog(t *testing.T) {
	var buffer bytes.Buffer
	expect := "Hello World!\n"
	logger := NewLogger(&buffer, &buffer)
	logger.Debug("Hello %s!", "World")

	if !strings.Contains(buffer.String(), expect) {
		t.Fatalf("Expected '%v' to contain '%v'", buffer.String(), expect)
	}

	logger.Debug("Another log %s", "entry")
	expect = "Another log entry\n"

	if !strings.Contains(buffer.String(), expect) {
		t.Fatalf("Expected '%v' to contain '%v'", buffer.String(), expect)
	}
}

func TestLoggerLevel(t *testing.T) {
	logger := NewLogger(nil, nil)
	logger.SetLevel(LevelDebug)

	if logger.level != LevelDebug || logger.GetLevel() != LevelDebug {
		t.Fatal("Unexpected log level")
	}

	logger.SetLevel(LevelInfo)

	if logger.level != LevelInfo || logger.GetLevel() != LevelInfo {
		t.Fatal("Unexpected log level")
	}

	logger.SetLevel(LevelWarning)

	if logger.level != LevelWarning || logger.GetLevel() != LevelWarning {
		t.Fatal("Unexpected log level")
	}

	logger.SetLevel(LevelError)

	if logger.level != LevelError || logger.GetLevel() != LevelError {
		t.Fatal("Unexpected log level")
	}
}

func TestLoggerFormatterAndOutput(t *testing.T) {
	var buffer bytes.Buffer
	expect := "Hello World!\n"
	logger := NewLogger(&buffer, &buffer)
	logger.Debug("Hello %s!", "World")

	if !strings.Contains(buffer.String(), expect) {
		t.Fatalf("Expected '%v' to contain '%v'", buffer.String(), expect)
	}

	formatter := NewDiscardFormatter()
	logger.SetFormatter(formatter)

	if logger.GetFormatter() != formatter {
		t.Fatal("Unexpected formatter")
	}

	var newBuffer bytes.Buffer
	logger.SetOut(LevelDebug, &newBuffer)
	logger.Debug("Another log %s", "entry")

	if newBuffer.String() != "" {
		t.Fatalf("Expected log to be empty but was: %v", buffer.String())
	}
}

func TestLoggerFatal(t *testing.T) {
	var stderr bytes.Buffer

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Fatal must panic")
		}

		if !strings.Contains(stderr.String(), "Fatal message") {
			t.Fatalf("Log must contain error message")
		}
	}()

	logger := NewLogger(nil, &stderr)
	logger.Fatal("Fatal %v", "message")
}

func TestNewLoggerRollingFileAppender(t *testing.T) {
	if err := os.RemoveAll("out"); err != nil {
		t.Fatal(err)
	}

	stdRFA, err := NewRollingFileAppender(0, 0, 0, "out", &testNameSchema{name: ".std"})

	if err != nil {
		t.Fatal(err)
	}

	errRFA, err := NewRollingFileAppender(0, 0, 0, "out", &testNameSchema{name: ".err"})

	if err != nil {
		t.Fatal(err)
	}

	logger := NewLogger(stdRFA, errRFA)
	logger.Info("info")
	logger.Error("error")

	if err := stdRFA.Close(); err != nil {
		t.Fatal(err)
	}

	if err := errRFA.Close(); err != nil {
		t.Fatal(err)
	}

	stdFile, err := os.ReadFile("out/1_log.std.txt")

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(stdFile), "info") {
		t.Fatalf("Info log must contain log output, but was: %v", string(stdFile))
	}

	errFile, err := os.ReadFile("out/1_log.err.txt")

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(errFile), "error") {
		t.Fatalf("Error log must contain log output, but was: %v", string(errFile))
	}
}
