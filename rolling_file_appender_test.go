package logbuch

import (
	"fmt"
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

func TestRollingFileAppender(t *testing.T) {
	if err := os.RemoveAll("out"); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir("out", 0777); err != nil {
		t.Fatal(err)
	}

	naming := &testNameSchema{}
	rfa, err := NewRollingFileAppender(3, 8, 10, "out", naming)

	if err != nil {
		t.Fatalf("Appender must be created, but was: %v", err)
	}

	for i := 0; i < 3; i++ {
		if n, err := rfa.Write([]byte("Hello World!\n")); err != nil || n != 13 {
			t.Fatalf("Log output must have been written, but was: %v %v", err, n)
		}
	}

	if err := rfa.Close(); err != nil {
		t.Fatalf("Appender must be closed, but was: %v", err)
	}
}
