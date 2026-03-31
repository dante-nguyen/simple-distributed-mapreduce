package errorsx

import "fmt"

func Wrap(message string, cause error) error {
	return fmt.Errorf("%s: %w", message, cause)
}

func WrapAsMessage(message string, cause error) string {
	return fmt.Sprintf("%s: %s", message, cause)
}
