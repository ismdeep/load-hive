package hive

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

var defaultConfigContent = `hive:
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
    worker: 16
  - ip: 10.20.30.43
    worker: 16
  - ip: 10.20.30.44
    worker: 16

extra-hosts:
  'example.com': 10.20.30.50
`

var LoadHiveProjectDir string

var h *Hive

func TestMain(m *testing.M) {
	var err error
	LoadHiveProjectDir, err = os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(fmt.Sprintf("%s/.load-hive", LoadHiveProjectDir), 0755); err != nil {
		panic(err)
	}

	if err := os.WriteFile(fmt.Sprintf("%s/.load-hive/default.yaml", LoadHiveProjectDir), []byte(defaultConfigContent), 0644); err != nil {
		panic(err)
	}

	tmpHive, err := New(LoadHiveProjectDir, "default")
	if err != nil {
		panic(err)
	}

	h = tmpHive

	raw, err := yaml.Marshal(h.config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("got = \n%v", string(raw))

	m.Run()

	_ = os.RemoveAll(LoadHiveProjectDir)
}
