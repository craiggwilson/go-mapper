package mapper

import (
	"reflect"
)

// TypeMapper handles mapping from src to dst.
type TypeMapper interface {
	// Dst is the type of the destination.
	Dst() reflect.Type
	// Src is the type of the source.
	Src() reflect.Type
	// Map performs the mapping to dst from src.
	Map(ctx Context, dst interface{}, src interface{}) error
}

// FunctionTypeMapper implements the TypeMapper interface by invoking a function
type FunctionTypeMapper struct {
	dst reflect.Type
	src reflect.Type
	f func(dst interface{}, src interface{}) error
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
	return tm.f(dst, src)
}
