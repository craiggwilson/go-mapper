package core

import (
	"reflect"
)

// Mapper handles mapping from src to dst.
type Mapper interface {
	// Dst is the type of the destination.
	Dst() reflect.Type
	// Src is the type of the source.
	Src() reflect.Type
	// Map performs the mapping to dst from src.
	Map(ctx Context, dst reflect.Value, src reflect.Value) error
}
