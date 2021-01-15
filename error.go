package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Error in recurerr is built with scope, or hierarchy in mind. This struct has
// a "parrent" error and its "children" errors: children errors belong to the
// parent error.
type Error struct {
	err error
	underlying []error
}

func New(err error, underlying...error) error {
	return &Error{
		err: err,
		underlying: underlying,
	}
}

func Newf(err error, format string, args...interface{}) error {
	return &Error{
		err: err,
		underlying: []error{fmt.Errorf(format, args...)},
	}
}

// Unwrap is for compatibility to the official errors.
func (e *Error) Unwrap() error {
	if len(e.underlying) > 0 {
		return e.underlying[0]
	}
	return nil
}

func (e *Error) Error() string {
	var sb strings.Builder
	sb.WriteString(e.err.Error())
	sb.WriteString(": ")
	if len(e.underlying) > 1 {
		sb.WriteByte('[')
	}
	for i, err := range e.underlying {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Error())
	}
	if len(e.underlying) > 1 {
		sb.WriteByte(']')
	}
	return sb.String()
}

func Is(s error, expect error) bool {
	if e, ok := s.(*Error); ok {
		if Is(e.err, expect) {
			return true
		}
		for _, err := range e.underlying {
			if Is(err, expect) {
				return true
			}
		}
	}
	return errors.Is(s, expect)
}

func As(s error, expect interface{}) bool {
	if e, ok := s.(*Error); ok {
		if As(e.err, expect) {
			return true
		}
		for _, err := range e.underlying {
			if As(err, expect) {
				return true
			}
		}
	}
	return errors.As(s, expect)
}
