package auto_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/craiggwilson/go-mapper/pkg/auto"
)

func TestSimple(t *testing.T) {
	t.Parallel()
	type customer struct {
		Name string
	}
	type customerDTO struct {
		Name string
	}

	ap := auto.NewProvider()
	ap.Add(
		reflect.TypeOf(new(customerDTO)),
		reflect.TypeOf(new(customer)),
	)

	mappers, err := ap.Mappers()
	require.NoError(t, err)

	src := customer{
		Name: "Blockus",
	}
	var dst customerDTO
	err = mappers[0].Map(nil, reflect.ValueOf(&dst), reflect.ValueOf(&src))
	require.NoError(t, err)
	require.Equal(t, src.Name, dst.Name)
}

func TestTypeConversion(t *testing.T) {
	t.Parallel()

	t.Run("int -> string", func(t *testing.T) {
		t.Parallel()
		type customer struct {
			Name int
		}
		type customerDTO struct {
			Name string
		}

		ap := auto.NewProvider()
		ap.Add(
			reflect.TypeOf(new(customerDTO)),
			reflect.TypeOf(new(customer)),
		)

		mappers, err := ap.Mappers()
		require.NoError(t, err)

		src := customer{
			Name: 42,
		}
		var dst customerDTO
		err = mappers[0].Map(nil, reflect.ValueOf(&dst), reflect.ValueOf(&src))
		require.NoError(t, err)
		require.Equal(t, strconv.Itoa(src.Name), dst.Name)
	})

	t.Run("string -> int", func(t *testing.T) {
		t.Parallel()
		type order struct {
			Number string
		}
		type orderDTO struct {
			Number int
		}

		ap := auto.NewProvider()
		ap.Add(
			reflect.TypeOf(new(orderDTO)),
			reflect.TypeOf(new(order)),
		)

		mappers, err := ap.Mappers()
		require.NoError(t, err)

		src := order{
			Number: "42",
		}
		var dst orderDTO
		err = mappers[0].Map(nil, reflect.ValueOf(&dst), reflect.ValueOf(&src))
		require.NoError(t, err)
		require.Equal(t, src.Number, strconv.Itoa(dst.Number))
	})

	t.Run("unsupported", func(t *testing.T) {
		t.Parallel()
		type order struct {
			Number string
		}
		type orderDTO struct {
			Number []int
		}

		ap := auto.NewProvider()
		ap.Add(
			reflect.TypeOf(new(orderDTO)),
			reflect.TypeOf(new(order)),
		)

		_, err := ap.Mappers()
		require.Error(t, err)
	})
}

func TestFieldOpts(t *testing.T) {
	t.Run("ignore a field with an unsupported type conversion", func(t *testing.T) {
		t.Parallel()
		type order struct {
			Number string
		}
		type orderDTO struct {
			Number []int
		}

		ap := auto.NewProvider()
		ap.Add(
			reflect.TypeOf(new(orderDTO)),
			reflect.TypeOf(new(order)),
			auto.WithStructField("Number", auto.WithFieldIgnore()),
		)

		_, err := ap.Mappers()
		require.NoError(t, err)
	})
}

func TestFlatten(t *testing.T) {
	t.Parallel()
	type customer struct {
		Name string
	}
	type order struct {
		Customer *customer
	}

	type orderDTO struct {
		CustomerName string
	}

	ap := auto.NewProvider()
	ap.Add(
		reflect.TypeOf(new(orderDTO)),
		reflect.TypeOf(new(order)),
	)

	mappers, err := ap.Mappers()
	require.NoError(t, err)

	src := order{
		Customer: &customer{
			Name: "Blockus",
		},
	}
	var dst orderDTO
	err = mappers[0].Map(nil, reflect.ValueOf(&dst), reflect.ValueOf(&src))
	require.NoError(t, err)
	require.Equal(t, src.Customer.Name, dst.CustomerName)
}
