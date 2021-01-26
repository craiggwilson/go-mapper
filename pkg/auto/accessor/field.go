package accessor

import (
	"reflect"
)

// Field returns an accessor for a struct field.
func Field(fld reflect.StructField) *FieldAccessor {
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
