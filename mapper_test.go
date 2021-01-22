package mapper_test

import (
	"strconv"
	"testing"

	"github.com/craiggwilson/go-mapper"
)

func TestStringToInt(t *testing.T) {
	c := mapper.NewConfig()
	c.AddTypeMapperFromFunc(func(dst *int, src string) error {
		i, err := strconv.ParseInt(src, 10, 32)
		if err != nil {
			return err
		}

		i32 := int(i)
		*dst = i32
		return nil
	})

	m := c.Build()

	var i int
	err := m.Map(&i, "42")
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if i != 42 {
		t.Fatalf("expected 42, but got %v", i)
	}
}