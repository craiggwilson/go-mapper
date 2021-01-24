package convert

import (
	"fmt"
	"reflect"
	"strconv"
)

// Converter implementations handle converting from one type to another.
type Converter interface {
	// Convert does a conversion from the src into the dst.
	Convert(dst reflect.Value, src reflect.Value) error
}

// ConverterFunc is a functional implementation of the Converter interface.
type ConverterFunc func(dst reflect.Value, src reflect.Value) error

// Convert implements the Converter interface.
func (f ConverterFunc) Convert(dst reflect.Value, src reflect.Value) error {
	return f(dst, src)
}

// Convert is the default implementation of the Converter interface.
func Convert(dst reflect.Value, src reflect.Value) error {
	switch dst.Kind() {
	case reflect.Int:
		return toInt(dst, src)
	default:
		return fmt.Errorf("no converter for %q available", dst.Type())
	}
}

func toInt(dst reflect.Value, src reflect.Value) error {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dst.SetInt(src.Int())
	case reflect.String:
		i64, err := strconv.ParseInt(src.String(), 10, 64)
		if err != nil {
			return err
		}

		dst.SetInt(i64)
	}

	return nil
}
