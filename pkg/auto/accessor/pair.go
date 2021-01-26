package accessor

import (
	"reflect"
)

// Pair makes an PairAccessor to chain accessors together.
func Pair(first, second Accessor) *PairAccessor {
	return &PairAccessor{first, second}
}

// PairAccessor chains access to a field through a pair of getters, where the results of the first
// are piped to the second.
type PairAccessor struct {
	first  Accessor
	second Accessor
}

// Name implements the Accessor interface.
func (a *PairAccessor) Name() string {
	return a.first.Name() + "." + a.second.Name()
}

// Type implements the Accessor interface.
func (a *PairAccessor) Type() reflect.Type {
	return a.second.Type()
}

// ValueFrom implements the Accessor interface.
func (a *PairAccessor) ValueFrom(v reflect.Value) reflect.Value {
	v = a.first.ValueFrom(v)
	if v.IsNil() {
		return v
	}

	return a.second.ValueFrom(v)
}