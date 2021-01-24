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
		Transaction string
	}

	type orderDTO struct {
		ID int
		CustomerName string
		Transaction int
	}

	ap := auto.NewProvider()
	ap.AddStruct(func(_ *orderDTO, _ *order) {})

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
