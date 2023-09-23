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
	return fmt.Sprintf("%s '%s'. %s", e.Message, e.FieldName, e.Inner)
}
