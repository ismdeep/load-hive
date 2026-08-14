package hive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLoadsConfigAndAppliesDefaults(t *testing.T) {
	h := newTestHive(t, defaultConfigContent)

	if h.configName != "default" {
		t.Fatalf("configName = %q, want default", h.configName)
	}
	if h.config.Hive.UUID != "a1" {
		t.Fatalf("uuid = %q, want a1", h.config.Hive.UUID)
	}
	if h.config.Hive.UserCount != 1000 {
		t.Fatalf("user_count = %d, want 1000", h.config.Hive.UserCount)
	}
	if h.config.Hive.SpawnRate != 100 {
		t.Fatalf("spawn_rate = %d, want 100", h.config.Hive.SpawnRate)
	}
	if h.config.Nodes[0].Worker != 16 {
		t.Fatalf("positive worker should be preserved, got %d", h.config.Nodes[0].Worker)
	}
	for i := 1; i < len(h.config.Nodes); i++ {
		if h.config.Nodes[i].Worker != 8 {
			t.Fatalf("nodes[%d].worker = %d, want 8", i, h.config.Nodes[i].Worker)
		}
	}
	if h.MainDir() != "/var/lib/lh-a1-main" {
		t.Fatalf("MainDir() = %q", h.MainDir())
	}
	if h.NodeDir() != "/var/lib/lh-a1-node" {
		t.Fatalf("NodeDir() = %q", h.NodeDir())
	}
}

func TestNewKeepsPositiveLoadSettings(t *testing.T) {
	h := newTestHive(t, `hive:
  uuid: hive2
  port: "8080"
  user_count: 42
  spawn_rate: 7
`)

	if h.config.Hive.UserCount != 42 {
		t.Fatalf("user_count = %d, want 42", h.config.Hive.UserCount)
	}
	if h.config.Hive.SpawnRate != 7 {
		t.Fatalf("spawn_rate = %d, want 7", h.config.Hive.SpawnRate)
	}
}

func TestNewReturnsErrorForMissingConfigFile(t *testing.T) {
	_, err := New(t.TempDir(), "missing")
	if err == nil {
		t.Fatal("New() error = nil, want missing file error")
	}
}

func TestNewReturnsErrorForInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, "default", "hive: [")

	_, err := New(dir, "default")
	if err == nil {
		t.Fatal("New() error = nil, want yaml error")
	}
}

func TestNewReturnsErrorForMissingRequiredFields(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantError string
	}{
		{
			name: "missing uuid",
			content: `hive:
  port: "8080"
`,
			wantError: "config hive uuid not set",
		},
		{
			name: "missing port",
			content: `hive:
  uuid: hive1
`,
			wantError: "config hive port not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeConfig(t, dir, "default", tt.content)

			_, err := New(dir, "default")
			if err == nil {
				t.Fatal("New() error = nil")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("New() error = %q, want contains %q", err.Error(), tt.wantError)
			}
		})
	}
}

func TestWriteConfigHelperWritesNamedConfig(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, "custom", defaultConfigContent)

	if _, err := os.Stat(filepath.Join(dir, ".load-hive", "custom.yaml")); err != nil {
		t.Fatalf("expected custom config file: %v", err)
	}
	if _, err := New(dir, "custom"); err != nil {
		t.Fatalf("New() with custom config error = %v", err)
	}
}
