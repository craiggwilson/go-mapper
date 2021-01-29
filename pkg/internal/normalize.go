package internal

import (
	"fmt"
	"reflect"
)

func EnsureSettableDst(dst reflect.Value) reflect.Value {
	if !dst.IsValid() {
		panic(fmt.Errorf("dst is not valid"))
	}
	if dst.Kind() != reflect.Ptr {
		panic(fmt.Errorf("dst must be a pointer, but was %v", dst.Kind()))
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
		panic(fmt.Errorf("dst could not be set"))
	}

	return dst
}

func UnwrapPtrValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

func UnwrapPtrType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}