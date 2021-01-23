package core_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/craiggwilson/go-mapper/pkg/core"
)

func TestFunctionTypeMapper(t *testing.T) {
	tm := core.MapperFromFunc(func(dst *int, src string) error {
		i, err := strconv.ParseInt(src, 10, 32)
		if err != nil {
			return err
		}

		i32 := int(i)
		*dst = i32
		return nil
	})

	var i int
	err := tm.Map(nil, reflect.ValueOf(&i), reflect.ValueOf("42"))
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if i != 42 {
		t.Fatalf("expected 42, but got %v", i)
	}
}
