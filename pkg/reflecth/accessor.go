package reflecth

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

// NewFieldAccessor makes a FieldAccessor.
func NewFieldAccessor(fld reflect.StructField) *FieldAccessor {
	return &FieldAccessor{fld}
}

// FieldAccessor retrieves a value from a field in a struct.
type FieldAccessor struct {
	fld reflect.StructField
}

// Name implements the Accessor interface.
func (a *FieldAccessor) Name() string {
	return a.fld.Name
}

// Type implements the Accessor interface.
func (a *FieldAccessor) Type() reflect.Type {
	return a.fld.Type
}

// ValueFrom implements the Accessor interface.
func (a *FieldAccessor) ValueFrom(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	return reflect.Indirect(v).FieldByIndex(a.fld.Index)
}

// NewAccessorPair makes an AccessorPair.
func NewAccessorPair(first, second Accessor) *AccessorPair {
	return &AccessorPair{first, second}
}

// AccessorPair chains access to a field through a pair of getters, where the results of the first
// are piped to the second.
type AccessorPair struct {
	first  Accessor
	second Accessor
}

// Name implements the Accessor interface.
func (a *AccessorPair) Name() string {
	return a.first.Name() + "." + a.second.Name()
}

// Type implements the Accessor interface.
func (a *AccessorPair) Type() reflect.Type {
	return a.second.Type()
}

// ValueFrom implements the Accessor interface.
func (a *AccessorPair) ValueFrom(v reflect.Value) reflect.Value {
	v = a.first.ValueFrom(v)
	if v.IsNil() {
		return v
	}

	return a.second.ValueFrom(v)
}
