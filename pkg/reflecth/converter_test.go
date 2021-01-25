package reflecth_test

import (
	"reflect"
	"testing"

	"github.com/craiggwilson/go-mapper/pkg/reflecth"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		src interface{}
		expected interface{}
	} {
		{"string to int", "10", 10},
		{"string to *int", "10", ptrTo(10)},
		{"string to **int", "10", ptrTo(ptrTo(10))},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			src := reflect.ValueOf(tc.src)
			dst := reflect.New(reflect.TypeOf(tc.expected))

			converter, err := reflecth.ConverterFor(dst.Type(), src.Type())
			if err != nil {
				t.Fatalf("expected no error, but got %v", err)
			}

			err = converter.Convert(dst, src)
			if err != nil {
				t.Fatalf("expected no error, but got %v", err)
			}

			if !reflect.DeepEqual(tc.expected, dst.Elem().Interface()) {
				t.Fatalf("expected %v, but got %v", tc.expected, dst.Elem().Interface())
			}
		})
	}
}

func ptrTo(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	p := reflect.New(v.Type())
	p.Elem().Set(v)
	return p.Interface()
}
