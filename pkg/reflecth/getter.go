package reflecth

import (
	"reflect"
)

// Getter retrieves a value.
type Getter interface {
	// In is the type of the input value.
	In() reflect.Type
	// Out is the type of the returned value.
	Out() reflect.Type
	// ValueFrom retrieves a value from the provided value.
	ValueFrom(v reflect.Value) reflect.Value
}

// NewFieldGetter makes a FieldGetter.
func NewFieldGetter(in reflect.Type, fld reflect.StructField) *FieldGetter {
	return &FieldGetter{in, fld}
}

// FieldGetter retrieves a value from a field in a struct.
type FieldGetter struct {
	in reflect.Type
	fld reflect.StructField
}

// In implements the Getter interface.
func (a *FieldGetter) In() reflect.Type {
	return a.in
}

// Out implements the Getter interface.
func (a *FieldGetter) Out() reflect.Type {
	return a.fld.Type
}

// ValueFrom implements the Getter interface.
func (a *FieldGetter) ValueFrom(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	return reflect.Indirect(v).FieldByIndex(a.fld.Index)
}

// NewAccessorPair makes an GetterPair.
func NewAccessorPair(first, second Getter) *GetterPair {
	return &GetterPair{first, second}
}

// GetterPair chains access to a field through a pair of getters, where the results of the first
// are piped to the second.
type GetterPair struct {
	first  Getter
	second Getter
}

// In implements the Getter interface.
func (a *GetterPair) In() reflect.Type {
	return a.first.In()
}

// Out implements the Getter interface.
func (a *GetterPair) Out() reflect.Type {
	return a.second.Out()
}

// ValueFrom implements the Getter interface.
func (a *GetterPair) ValueFrom(v reflect.Value) reflect.Value {
	v = a.first.ValueFrom(v)
	if v.IsNil() {
		return v
	}

	return a.second.ValueFrom(v)
}
