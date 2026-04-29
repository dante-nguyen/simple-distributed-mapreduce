package logx

import "log"

func Err(ctx string, err error) {
	log.Printf("[ERROR] %s: %s", ctx, err)
}
