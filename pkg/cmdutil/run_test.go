package cmdutil

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunUsesDirAndEnv(t *testing.T) {
	dir := t.TempDir()
	realDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("filepath.EvalSymlinks() error = %v", err)
	}

	output, err := Run(dir, map[string]string{"LH_TEST_VALUE": "ok"}, "sh", "-c", "printf '%s:%s' \"$PWD\" \"$LH_TEST_VALUE\"")
	if err != nil {
		t.Fatalf("Run() error = %v, output = %q", err, output)
	}
	if output != realDir+":ok" {
		t.Fatalf("Run() output = %q, want %q", output, realDir+":ok")
	}
}

func TestRunReturnsCombinedOutputOnFailure(t *testing.T) {
	output, err := Run(t.TempDir(), nil, "sh", "-c", "printf stdout; printf stderr >&2; exit 7")
	if err == nil {
		t.Fatal("Run() error = nil, want exit error")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("Run() error type = %T, want *exec.ExitError", err)
	}
	if !strings.Contains(output, "stdout") || !strings.Contains(output, "stderr") {
		t.Fatalf("Run() output = %q, want stdout and stderr", output)
	}
}
