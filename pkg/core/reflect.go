package core

import (
	"reflect"
)

var (
	tContext        = reflect.TypeOf((*Context)(nil)).Elem()
	tErr            = reflect.TypeOf((*error)(nil)).Elem()
)