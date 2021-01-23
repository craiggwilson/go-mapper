package mapper_test

import (
	"strconv"
	"testing"

	"github.com/craiggwilson/go-mapper"
)

func TestStringToInt(t *testing.T) {
	c := mapper.NewStaticConfig()
	c.Add(mapper.NewFunctionTypeMapper(func(dst *int, src string) error {
		i, err := strconv.ParseInt(src, 10, 32)
		if err != nil {
			return err
		}

		i32 := int(i)
		*dst = i32
		return nil
	}))

	m, err := mapper.New(c)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	var i int
	err = m.Map(&i, "42")
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if i != 42 {
		t.Fatalf("expected 42, but got %v", i)
	}
}

func TestStructToStruct(t *testing.T) {
	type customer struct {
		Name string
	}
	type order struct {
		ID int
		Customer *customer
	}

	type orderDTO struct {
		ID int
		CustomerName string
	}

	c := mapper.NewAutoConfig()

	c.Add(func(_ *orderDTO, _ *order, opts *mapper.AutoTypeOptions) {
		opts.Field("CustomerName", func(dst *string, src *order) error {
			*dst = src.Customer.Name
			return nil
		})
	})

	m, err := mapper.New(c)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	src := order{
		ID: 10,
		Customer: &customer {
			Name: "Blockus",
		},
	}
	var dst orderDTO
	err = m.Map(&dst, &src)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	//if dst.ID != src.ID {
	//	t.Fatalf("expected %d, but got %d", src.ID, dst.ID)
	//}
	if dst.CustomerName != src.Customer.Name {
		t.Fatalf("expected %q, but got %q", src.Customer.Name, dst.CustomerName)
	}
}
