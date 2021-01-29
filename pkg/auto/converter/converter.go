package converter

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/craiggwilson/go-mapper/pkg/internal"
)

// Factory returns a Converter.
type Factory interface {
	// For returns a converter for the src and dest.
	ConverterFor(dst reflect.Type, src reflect.Type) (Converter, error)
}

// FactoryFunc is a functional implementation of a Factory.
type FactoryFunc func(dst reflect.Type, src reflect.Type) (Converter, error)

// Convert implements the Converter interface.
func (f FactoryFunc) ConverterFor(dst reflect.Type, src reflect.Type) (Converter, error) {
	return f(dst, src)
}

// Converter implementations handle converting from one type to another.
type Converter interface {
	// Convert does a conversion from the src into the dst.
	Convert(dst reflect.Value, src reflect.Value) error
}

// Func is a function implementation of a Converter.
type Func func(dst reflect.Value, src reflect.Value) error

// Convert implements the Converter interface.
func (f Func) Convert(dst reflect.Value, src reflect.Value) error {
	return f(dst, src)
}

// For is the default implementation of a Factory.
func For(dst reflect.Type, src reflect.Type) (Converter, error) {
	dst = internal.UnwrapPtrType(dst)
	src = internal.UnwrapPtrType(src)
	if dst.AssignableTo(src) {
		return nil, nil
	}
	switch dst.Kind() {
	case reflect.Int:
		return toIntConverter(src)
	default:
		return nil, fmt.Errorf("cannot convert from %v to %v", src, dst)
	}
}

func toIntConverter(src reflect.Type) (Converter, error) {
	switch src.Kind() {
	case reflect.Int:
		return nil, nil
	case reflect.String:
		return Func(stringToInt), nil
	default:
		return nil, fmt.Errorf("no converter for string -> %v available", src)
	}
}

func stringToInt(dst reflect.Value, src reflect.Value) error {
	dst = internal.EnsureSettableDst(dst)
	src = internal.UnwrapPtrValue(src)

	if src.IsZero() {
		return nil
	}

	i, err := strconv.Atoi(src.String())
	if err != nil {
		return err
	}

	dst.SetInt(int64(i))
	return nil
}
