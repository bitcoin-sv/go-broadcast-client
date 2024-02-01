package outter_errors

import (
	"fmt"
	"strings"
)

func New(msg string, inner error) error {
	return &outterErr{
		msg:   msg,
		inner: inner,
	}
}

type outterErr struct {
	msg   string
	inner error
}

func (e *outterErr) Error() string {
	return e.error(1)
}

func (e *outterErr) Unwrap() error {
	return e.inner
}

func (e *outterErr) error(level uint8) string {

	if e.inner != nil {
		var innerMsg string
		if oie, ok := e.inner.(*outterErr); ok {
			innerMsg = oie.error(level + 1)
		} else {
			innerMsg = e.inner.Error()
		}

		indent := strings.Repeat("\t", int(level))
		return fmt.Sprintf("%s\n%sinner error: %s", e.msg, indent, innerMsg)

	} else {
		return e.msg
	}
}
