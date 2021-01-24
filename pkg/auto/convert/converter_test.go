package convert_test

import (
	"reflect"
	"testing"

	"github.com/craiggwilson/go-mapper/pkg/auto/convert"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		src interface{}
		expected interface{}
	} {
		{"string to int", "10", 10},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			src := reflect.ValueOf(tc.src)
			dst := reflect.New(reflect.TypeOf(tc.expected))

			err := convert.Convert(dst.Elem(), src)
			if err != nil {
				t.Fatalf("expected no error, but got %v", err)
			}

			if !reflect.DeepEqual(tc.expected, dst.Elem().Interface()) {
				t.Fatalf("expected %v, but got %v", tc.expected, dst.Elem().Interface())
			}
		})
	}
}
