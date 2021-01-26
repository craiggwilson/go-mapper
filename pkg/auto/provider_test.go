package auto_test

import (
	"reflect"
	"strconv"
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
	ap.Add(
		reflect.TypeOf(new(orderDTO)),
		reflect.TypeOf(new(order)),
	)

	src := order{
		ID: 10,
		Customer: &customer {
			Name: "Blockus",
		},
		Transaction: "42",
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
		t.Fatalf("expected %s, but got %s", src.Customer.Name, dst.CustomerName)
	}
	if trans, _ := strconv.Atoi(src.Transaction); dst.Transaction != trans {
		t.Fatalf("expected %d, but got %d", trans, dst.Transaction)
	}
}
