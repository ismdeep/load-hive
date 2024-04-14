package hive

import "testing"

func TestHiveDown(t *testing.T) {
	PanicIf(h.Down())
}
