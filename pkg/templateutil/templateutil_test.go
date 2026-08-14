package templateutil

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	got, err := Generate("hello {{ .Name }}", map[string]string{"Name": "load-hive"})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got != "hello load-hive" {
		t.Fatalf("Generate() = %q", got)
	}
}

func TestGenerateReturnsParseError(t *testing.T) {
	_, err := Generate("hello {{", nil)
	if err == nil {
		t.Fatal("Generate() error = nil, want parse error")
	}
}

func TestGenerateReturnsExecuteError(t *testing.T) {
	_, err := Generate("hello {{ .Name }}", struct{}{})
	if err == nil {
		t.Fatal("Generate() error = nil, want execute error")
	}
	if !strings.Contains(err.Error(), "Name") {
		t.Fatalf("Generate() error = %q, want field name in error", err.Error())
	}
}
