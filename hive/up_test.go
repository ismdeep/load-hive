package hive

import "testing"

func TestHiveUp(t *testing.T) {
	PanicIf(h.Up())
}
