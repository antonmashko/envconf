package envconf

import (
	"errors"
	"fmt"
)

var (
	// ErrNilData mean that exists nil pointer inside data struct
	ErrNilData               = errors.New("nil data")
	ErrUnsupportedType       = errors.New("unsupported type")
	ErrConfigurationNotFound = errors.New("configuration not found")
)

type Error struct {
	Inner     error
	Message   string
	FieldName string
}

func (e *Error) Error() string {
	msg := e.Message
	if e.FieldName != "" {
		msg = fmt.Sprintf("%s: %s", e.FieldName, msg)
	}
	if e.Inner != nil {
		msg += " " + e.Inner.Error()
	}
	return msg
}

func (e *Error) Unwrap() error {
	return e.Inner
}
