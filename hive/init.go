package hive

import (
	"fmt"
	"time"
)

func Log(level string, speaker string, format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	fmt.Printf("%s [ %s ] %s: %v\n", time.Now().Format(time.RFC3339), level, speaker, s)
}
