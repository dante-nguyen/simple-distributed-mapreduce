package errx

import "fmt"

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
