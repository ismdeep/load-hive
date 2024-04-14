package hive

import "fmt"

func Log(level string, speaker string, format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s: %v\n", level, speaker, s)
}
