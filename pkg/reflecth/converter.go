package reflecth

import (
	"fmt"
	"reflect"
	"strconv"
)

// ConverterFactory returns a Converter.
type ConverterFactory interface {
	// ConverterFor returns a converter for the src and dest.
	ConverterFor(dst reflect.Type, src reflect.Type) (Converter, error)
}

// ConverterFactoryFunc is a functional implementation of a ConverterFactory.
type ConverterFactoryFunc func(dst reflect.Type, src reflect.Type) (Converter, error)

// Convert implements the Converter interface.
func (f ConverterFactoryFunc) ConverterFor(dst reflect.Type, src reflect.Type) (Converter, error) {
	return f(dst, src)
}

// Converter implementations handle converting from one type to another.
type Converter interface {
	// Convert does a conversion from the src into the dst.
	Convert(dst reflect.Value, src reflect.Value) error
}

// ConverterFunc is a function implementation of a Converter.
type ConverterFunc func(dst reflect.Value, src reflect.Value) error

// Convert implements the Converter interface.
func (f ConverterFunc) Convert(dst reflect.Value, src reflect.Value) error {
	return f(dst, src)
}

// ConverterFor is the default implementation of a ConverterFactory.
func ConverterFor(dst reflect.Type, src reflect.Type) (Converter, error) {
	dst = unwrapPtr(dst)
	src = unwrapPtr(src)
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
		return ConverterFunc(stringToInt), nil
	default:
		return nil, fmt.Errorf("no converter for string -> %v available", src)
	}
}

func stringToInt(dst reflect.Value, src reflect.Value) error {
	dst, err := findSetter(dst)
	if err != nil {
		return err
	}

	i64, err := strconv.ParseInt(src.String(), 10, 64)
	if err != nil {
		return err
	}

	dst.SetInt(i64)
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
