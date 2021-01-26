package accessor

import (
	"reflect"
)

// Accessor retrieves a value.
type Accessor interface {
	// Name indicates the source of the value.
	Name() string
	// Type is the type of the returned value.
	Type() reflect.Type
	// ValueFrom retrieves a value from the provided value.
	ValueFrom(v reflect.Value) reflect.Value
}
