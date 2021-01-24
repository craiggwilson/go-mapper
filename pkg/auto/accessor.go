package auto

import (
	"reflect"
)

// Accessor retrieves a value from another value.
type Accessor interface {
	// In is the type of the input value.
	In() reflect.Type
	// Out is the type of the returned value.
	Out() reflect.Type
	// ValueFrom retrieves a value from the provided value.
	ValueFrom(v reflect.Value) reflect.Value
}

// NewFieldAccessor makes a FieldAccessor.
func NewFieldAccessor(in reflect.Type, fld reflect.StructField) *FieldAccessor {
	return &FieldAccessor{in, fld}
}

// FieldAccessor retrieves a value from a field in a struct.
type FieldAccessor struct {
	in reflect.Type
	fld reflect.StructField
}

// In implements the Accessor interface.
func (a *FieldAccessor) In() reflect.Type {
	return a.in
}

// Out implements the Accessor interface.
func (a *FieldAccessor) Out() reflect.Type {
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

// AccessorPair chains access to a field through a pair of accessors, where the results of the first
// are piped to the second.
type AccessorPair struct {
	first Accessor
	second Accessor
}

// In implements the Accessor interface.
func (a *AccessorPair) In() reflect.Type {
	return a.first.In()
}

// Out implements the Accessor interface.
func (a *AccessorPair) Out() reflect.Type {
	return a.second.Out()
}

// ValueFrom implements the Accessor interface.
func (a *AccessorPair) ValueFrom(v reflect.Value) reflect.Value {
	v = a.first.ValueFrom(v)
	if v.IsNil() {
		return v
	}

	return a.second.ValueFrom(v)
}