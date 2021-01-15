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
	err        error
	message    string
	underlying []error
}

// New accepts the parent error, following optional children errors,
// following optional format strings. Remaining arguments are discarded.
// That is <err> [(err1, err2, ..)] [(format, args)].
func New(nerr interface{}, args ...interface{}) error {
	err, ok := nerr.(error)
	if !ok {
		err = fmt.Errorf("%v", nerr)
	}
	if len(args) == 0 {
		return err
	}
	underlying := []error{}
	idx := len(args)
	for i, arg := range args {
		if e, ok := arg.(error); !ok {
			idx = i
			break
		} else {
			underlying = append(underlying, e)
		}
	}
	ret := &Error{
		err,
		"",
		underlying,
	}
	if idx < len(args) {
		if fmtstr, ok := args[idx].(string); ok {
			ret.message = fmt.Sprintf(fmtstr, args[idx+1:]...)
		}
	}
	return ret
}

// Message returns comments on the parent error.
func (e *Error) Message() string {
	return e.message
}

// Unwrap is for compatibility to the official errors.
func (e *Error) Unwrap() error {
	if len(e.underlying) > 0 {
		return e.underlying[0]
	}
	return nil
}

// Error implements error.
func (e *Error) Error() string {
	var sb strings.Builder
	sb.WriteString(e.err.Error())
	sb.WriteString(": ")
	sb.WriteString(e.message)
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
