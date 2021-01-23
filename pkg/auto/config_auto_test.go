package auto_test

import (
	"reflect"
	"testing"

	"github.com/craiggwilson/go-mapper/pkg/auto"
)

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

	ap := auto.NewProvider()
	ap.AddStruct(func(_ *orderDTO, _ *order, opts *auto.StructOptions) {
		opts.Field("CustomerName", func(dst *string, src *order) error {
			*dst = src.Customer.Name
			return nil
		})
	})

	src := order{
		ID: 10,
		Customer: &customer {
			Name: "Blockus",
		},
	}
	var dst orderDTO
	err := ap.Mappers()[0].Map(nil, reflect.ValueOf(&dst), reflect.ValueOf(&src))
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if dst.ID != src.ID {
		t.Fatalf("expected %d, but got %d", src.ID, dst.ID)
	}
	if dst.CustomerName != src.Customer.Name {
		t.Fatalf("expected %q, but got %q", src.Customer.Name, dst.CustomerName)
	}
}