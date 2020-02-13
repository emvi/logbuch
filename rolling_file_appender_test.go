package logbuch

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type testNameSchema struct {
	counter int
}

func (schema *testNameSchema) Name() string {
	schema.counter++
	return fmt.Sprintf("%d_log.txt", schema.counter)
}

func TestNewRollingFileAppender(t *testing.T) {
	if err := os.RemoveAll("out"); err != nil {
		t.Fatal(err)
	}

	rfa, err := NewRollingFileAppender(-1, -1, -1, "out", nil)

	if rfa != nil || err == nil || err.Error() != "filename schema must be specified" {
		t.Fatalf("Filename schema must be passed, but was: %v", err)
	}

	rfa, err = NewRollingFileAppender(-1, -1, -1, "out", &testNameSchema{})

	if err != nil {
		t.Fatalf("Rolling file appender must have been created, but was: %v", err)
	}

	if rfa.Files != defaultFiles || rfa.FileSize != defaultFileSize || rfa.maxBufferSize != defaultBufferSize {
		t.Fatalf("Default values not as expected: %v", rfa)
	}

	if _, err := os.Stat("out/1_log.txt"); os.IsNotExist(err) {
		t.Fatalf("First log file must exist, but was: %v", err)
	}

	if err := rfa.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRollingFileAppender_Write(t *testing.T) {
	if err := os.RemoveAll("out"); err != nil {
		t.Fatal(err)
	}

	rfa, err := NewRollingFileAppender(3, 5, 10, "out", &testNameSchema{})

	if err != nil {
		t.Fatalf("Appender must be created, but was: %v", err)
	}

	for i := 0; i < 4; i++ {
		if n, err := rfa.Write([]byte(fmt.Sprintf("%d234\n", i))); err != nil || n != 5 {
			t.Fatalf("Log output must have been written, but was: %v %v", err, n)
		}
	}

	if err := rfa.Close(); err != nil {
		t.Fatalf("Appender must be closed, but was: %v", err)
	}

	dir, err := ioutil.ReadDir("out")

	if err != nil {
		t.Fatal(err)
	}

	if len(dir) != 3 {
		t.Fatalf("Three files must have been created, but was: %v", len(dir))
	}

	if dir[0].Name() != "2_log.txt" || dir[1].Name() != "3_log.txt" || dir[2].Name() != "4_log.txt" {
		t.Fatalf("Files must have appropriate names, but was: %v %v %v", dir[0].Name(), dir[1].Name(), dir[2].Name())
	}

	if _, err := os.Stat("out/1_log.txt"); !os.IsNotExist(err) {
		t.Fatalf("First log file must not exist, but was: %v", err)
	}
}

func TestRollingFileAppender_Flush(t *testing.T) {
	if err := os.RemoveAll("out"); err != nil {
		t.Fatal(err)
	}

	rfa, err := NewRollingFileAppender(3, 100, 60, "out", &testNameSchema{})

	if err != nil {
		t.Fatalf("Appender must be created, but was: %v", err)
	}

	if _, err := rfa.Write([]byte("This fits into the buffer!\n")); err != nil {
		t.Fatalf("Log output must have been written, but was: %v", err)
	}

	content, err := ioutil.ReadFile("out/1_log.txt")

	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "" {
		t.Fatalf("File must be empty, but was: %v", string(content))
	}

	if err := rfa.Flush(); err != nil {
		t.Fatalf("Log must have been flushed, but was: %v", err)
	}

	content, err = ioutil.ReadFile("out/1_log.txt")

	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "This fits into the buffer!\n" {
		t.Fatalf("The log must have been saved, but was: %v", string(content))
	}

	if err := rfa.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRollingFileAppender_Close(t *testing.T) {
	if err := os.RemoveAll("out"); err != nil {
		t.Fatal(err)
	}

	rfa, err := NewRollingFileAppender(3, 100, 60, "out", &testNameSchema{})

	if err != nil {
		t.Fatalf("Appender must be created, but was: %v", err)
	}

	if _, err := rfa.Write([]byte("This fits into the buffer!\n")); err != nil {
		t.Fatalf("Log output must have been written, but was: %v", err)
	}

	if err := rfa.Close(); err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadFile("out/1_log.txt")

	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "This fits into the buffer!\n" {
		t.Fatalf("The log must have been saved, but was: %v", string(content))
	}
}
