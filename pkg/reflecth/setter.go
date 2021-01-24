package reflecth

import (
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/auto/convert"
)

// Setter sets a value.
type Setter interface {
	Type() reflect.Type
	Set(reflect.Value)
}

type Assigner interface {
	Assign(dst Setter, src Getter) error
}

type TypeConversionAssigner struct {
	converter convert.Converter
}


