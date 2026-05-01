package errx

import (
	"errors"
	"fmt"
)

func Chain(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var ret error = nil
	for i := len(errs) - 1; i >= 0; i-- {
		if ret == nil {
			ret = errs[i]
		} else {
			ret = fmt.Errorf("%w: %w", errs[i], ret)
		}
	}
	return ret
}

func OneOf(err error, errs ...error) bool {
	for _, target := range errs {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func WithContext(err error, ctx error) error {
	return Chain(ctx, err)
}

func WithContextMsg(err error, ctx string) error {
	return fmt.Errorf("%s: %w", ctx, err)
}
