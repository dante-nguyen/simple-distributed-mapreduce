package logx

import "log"

func Err(err error) {
	log.Printf("[ERROR] %s", err)
}
