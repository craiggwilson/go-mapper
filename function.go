package mapper

import (
	"fmt"
	"reflect"
)

// NewFunctionTypeMapper makes a FunctionTypeMapper. The fn argument must match the signature
// func(dst <type>, src <type>) error or func(ctx Context, dst <type>, src <type>). If fn is not a function,
// or it's signature does not match the requirements, a panic is raised.
func NewFunctionTypeMapper(fn interface{}) *FunctionTypeMapper {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("fn argument must be a func but got a %q", t.Kind()))
	}

	switch t.NumOut() {
	case 1:
		if !t.Out(0).AssignableTo(tErr) {
			panic(fmt.Sprintf("fn function must return an error, but returns %q", t.Out(0)))
		}
	default:
		panic(fmt.Errorf("fn function must return 1 value, but had %d", t.NumOut()))
	}

	argPos := 0
	switch t.NumIn() {
	case 3:
		if !t.In(0).AssignableTo(tContext) {
			panic(fmt.Errorf("fn function with 3 arguments must have a Context as the first, but got %q", t.In(0)))
		}
		argPos = 1
	case 2:
	default:
		panic(fmt.Errorf("fn function must have 2 or 3 arguments, but had %d", t.NumIn()))
	}

	v := reflect.ValueOf(fn)
	fnWrapper := func(ctx Context, dst interface{}, src interface{}) error {
		in := make([]reflect.Value, t.NumIn())
		if len(in) == 3 {
			in[0] = reflect.ValueOf(ctx)
		}
		in[argPos] = reflect.ValueOf(dst)
		in[argPos+1] = reflect.ValueOf(src)
		result := v.Call(in)

		if result[0].IsNil() {
			return nil
		}

		return result[0].Interface().(error)
	}

	return &FunctionTypeMapper{
		dst:   t.In(argPos),
		src:   t.In(argPos +1),
		mapFn: fnWrapper,
	}
}

// FunctionTypeMapper implements the TypeMapper interface by invoking a function
type FunctionTypeMapper struct {
	dst   reflect.Type
	src   reflect.Type
	mapFn func(ctx Context, dst interface{}, src interface{}) error
}

// Dst implements the TypeMapper interface.
func (tm *FunctionTypeMapper) Dst() reflect.Type {
	return tm.dst
}

// Src implements the TypeMapper interface.
func (tm *FunctionTypeMapper) Src() reflect.Type {
	return tm.src
}

// Map implements the TypeMapper interface.
func (tm *FunctionTypeMapper) Map(ctx Context, dst interface{}, src interface{}) error {
	return tm.mapFn(ctx, dst, src)
}
