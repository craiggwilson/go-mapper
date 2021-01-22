package mapper

import (
	"fmt"
	"reflect"
)

var (
	tContext = reflect.TypeOf((*Context)(nil)).Elem()
	tErr = reflect.TypeOf((*error)(nil)).Elem()
)

// NewConfig makes a Config.
func NewConfig() *Config {
	return &Config{}
}

// Config is used to build up a Mapper.
type Config struct {
	static staticTypeMapperProvider
	providers []TypeMapperProvider
}

// AddTypeMapper adds a TypeMapper to the static mapping list.
func (c *Config) AddTypeMapper(tm TypeMapper) error {
	return c.static.addTypeMapper(tm)
}

// AddTypeMapperFunc adds TypeMapper to the static mapping list using the provider function. The fn argument
// must match the signature func(dst <type>, src <type>) error or func(ctx Context, dst <type>, src <type>). If fn is not a function, or it's signature does
// not match the requirements, AddTypeMapperFromFunc will panic.
func (c *Config) AddTypeMapperFromFunc(fn interface{}) error {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("fn argument must be a func but got a %q", t.Kind()))
	}

	switch t.NumOut() {
	case 1:
		if !t.Out(0).AssignableTo(tErr) {
			panic(fmt.Sprintf("fn function must return an error, but returns %q", t.Out(0)))
		}
	default:
		if t.NumOut() != 1 {
			panic(fmt.Sprintf("fn function must 1 return value, but had %d", t.NumOut()))
		}
	}

	argPos := 0
	switch t.NumIn() {
	case 3:
		if !t.In(0).AssignableTo(tContext) {
			panic(fmt.Sprintf("fn function with 3 arguments first argument must be Context, but got %q", t.In(0)))
		}
		argPos = 1
	case 2:
	default:
		panic(fmt.Sprintf("fn function must have 2 or 3 arguments, but had %d", t.NumIn()))
	}

	v := reflect.ValueOf(fn)
	fnWrapper := func(dst interface{}, src interface{}) error {
		result := v.Call([]reflect.Value{
			reflect.ValueOf(dst),
			reflect.ValueOf(src),
		})

		if result[0].IsNil() {
			return nil
		}

		return result[0].Interface().(error)
	}

	return c.AddTypeMapper(&FunctionTypeMapper{
		dst: t.In(argPos),
		src: t.In(argPos +1),
		f: fnWrapper,
	})
}

func (c *Config) AddTypeMapperProvider(tmp TypeMapperProvider) {
	c.providers = append(c.providers, tmp)
}

func (c *Config) Build() *Mapper {
	return &Mapper{
		providers: append([]TypeMapperProvider{&c.static}, c.providers...),
	}
}
