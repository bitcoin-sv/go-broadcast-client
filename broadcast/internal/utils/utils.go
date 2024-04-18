package utils

import (
	"errors"
	"fmt"
)

func WithCause(err error, cause error) error {
	return errors.Join(err, fmt.Errorf("\tcaused by: %w\t", cause))
}
