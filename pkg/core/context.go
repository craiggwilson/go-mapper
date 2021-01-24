package core

import (
	"reflect"
)

type Context interface {
	Map(dst reflect.Value, src reflect.Value) error
}
