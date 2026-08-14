package hive

import (
	"os"
	"path/filepath"
	"testing"
)

const defaultConfigContent = `hive:
  uuid: a1
  port: 28373
  internal-api-port: 5555
  target: https://example.com

main:
  ip: 10.20.30.40

nodes:
  - ip: 10.20.30.41
    worker: 16
  - ip: 10.20.30.42
    worker: 0
  - ip: 10.20.30.43
    worker: -1

extra-hosts:
  'example.com': 10.20.30.50
`

func writeConfig(t *testing.T, dir, name, content string) {
	t.Helper()

	configDir := filepath.Join(dir, ".load-hive")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, name+".yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
}

func newTestHive(t *testing.T, content string) *Hive {
	t.Helper()

	dir := t.TempDir()
	writeConfig(t, dir, "default", content)

	h, err := New(dir, "default")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return h
}
