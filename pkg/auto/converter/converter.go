package converter

import (
	"fmt"
	"reflect"
	"strconv"
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
	dst = unwrapPtr(dst)
	src = unwrapPtr(src)
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
	dst, err := findSetter(dst)
	if err != nil {
		return err
	}

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

func findSetter(dst reflect.Value) (reflect.Value, error) {
	if dst.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("dst must be a pointer")
	}

	dst = dst.Elem()
	for dst.Kind() == reflect.Ptr {
		if !dst.Elem().IsValid() || dst.Elem().IsNil() {
			nv := reflect.New(dst.Type().Elem())
			dst.Set(nv)
		}

		dst = dst.Elem()
	}

	if !dst.CanSet() {
		return reflect.Value{}, fmt.Errorf("dst must be settable")
	}

	return dst, nil
}

func unwrapPtr(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
