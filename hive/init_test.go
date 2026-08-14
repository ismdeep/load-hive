package hive

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	output := captureStdout(t, func() {
		Log("INFO", "load-hive", "test %v", 1)
	})

	if !strings.Contains(output, "[ INFO ] load-hive: test 1") {
		t.Fatalf("Log() output = %q", output)
	}
	if ok := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T`).MatchString(output); !ok {
		t.Fatalf("Log() output should start with RFC3339-like timestamp, got %q", output)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = writePipe

	fn()

	if err := writePipe.Close(); err != nil {
		t.Fatalf("failed to close stdout pipe writer: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, readPipe); err != nil {
		t.Fatalf("failed to read stdout pipe: %v", err)
	}
	if err := readPipe.Close(); err != nil {
		t.Fatalf("failed to close stdout pipe reader: %v", err)
	}
	return buf.String()
}
