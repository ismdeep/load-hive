package hive

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

const (
	LoadHiveProjectDir = "/Users/ismdeep/Projects/github.com_ismdeep_load-hive-demo/main/load-hive-demo"
)

var h *Hive

func TestMain(m *testing.M) {
	tmpHive, err := New(LoadHiveProjectDir, "default")
	PanicIf(err)

	h = tmpHive

	raw, err := yaml.Marshal(h.config)
	PanicIf(err)

	fmt.Printf("got = \n%v", string(raw))

	m.Run()
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
