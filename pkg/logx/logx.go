package logx

import (
	"fmt"
	"log"
)

func Warn(msg string) {
	log.Printf("[WARNING] %s", msg)
}

func Warnf(format string, args ...any) {
	log.Printf("[WARNING] %s", fmt.Sprintf(format, args...))
}

func Err(err error) {
	log.Printf("[ERROR] %s", err)
}
