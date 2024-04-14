package hive

import "testing"

func TestLog(t *testing.T) {
	Log("INFO", "load-hive", "test %v", 1)
}
