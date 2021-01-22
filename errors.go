package mapper

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNoTypeMapperFound = errors.New("no type mapper found")

// DuplicateTypeError represents a configuration error that occurs when attempting to create a type mapper
// for a combination that already exists.
type DuplicateTypeError struct {
	Dst reflect.Type
	Src reflect.Type
}

// Error implements the error interface.
func (e *DuplicateTypeError) Error() string {
	return fmt.Sprintf("type map already exists for %q to %q", e.Src, e.Dst)
}

func newDuplicateTypeMapError(dst reflect.Type, src reflect.Type) *DuplicateTypeError {
	return &DuplicateTypeError{dst, src}
}
