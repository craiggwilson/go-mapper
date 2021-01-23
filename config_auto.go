package mapper

import (
	"fmt"
	"reflect"
)

// NewAutoConfig makes an AutoConfig.
func NewAutoConfig() *AutoConfig {
	return &AutoConfig{}
}

// AutoConfig is used to automatically map types following prescribed strategies for naming and type conversion.
type AutoConfig struct {
	// strategies

	tms []TypeMapper
}

// TypeMappers implements the TypeMapperProvider interface.
func (c *AutoConfig) TypeMappers() []TypeMapper {
	return c.tms
}

// Add registers a mapping. The fn argument must match the signature
// func(dst <type>, src <type>) or func(dst <type>, src <type>, cfg *AutoTypeOptions). If fn is not a function,
// or it's signature does not match the requirements, a panic is raised.
func (c *AutoConfig) Add(fn interface{}) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("fn argument must be a func but got a %q", t.Kind()))
	}

	switch t.NumOut() {
	case 0:
	default:
		panic(fmt.Errorf("fn function must have no return values, but had %d", t.NumOut()))
	}

	var opts AutoTypeOptions
	switch t.NumIn() {
	case 3:
		if !t.In(2).AssignableTo(tAutoTypeConfig) {
			panic(fmt.Errorf("fn function with 3 arguments must have *AutoTypeOptions as the last, but got %q", t.In(2)))
		}

		opts.dst = t.In(0)
		if opts.dst.Kind() != reflect.Ptr || opts.dst.Elem().Kind() != reflect.Struct {
			panic(fmt.Errorf("fn function's first argument must be a pointer to a struct"))
		}
		opts.src = t.In(1)

		v := reflect.ValueOf(fn)

		_ = v.Call([]reflect.Value{
			reflect.Zero(opts.dst),
			reflect.Zero(opts.src),
			reflect.ValueOf(&opts),
		})

	case 2:
		opts.dst = t.In(0)
		opts.src = t.In(1)
	default:
		panic(fmt.Errorf("fn function must have 2 or 3 arguments, but had %d", t.NumIn()))
	}

	tm := c.createTypeMapper(&opts)
	c.tms = append(c.tms, tm)
}

func (c *AutoConfig) createTypeMapper(opts *AutoTypeOptions) TypeMapper {
	return &FunctionTypeMapper{
		dst: opts.dst,
		src: opts.src,
		mapFn: func(ctx Context, dst interface{}, src interface{}) error {
			vDst := reflect.ValueOf(dst)
			vSrc := reflect.ValueOf(src)
			for _, fldopt := range opts.fldOpts {
				fv := reflect.Indirect(vDst).FieldByName(fldopt.dst.Name)
				if err := fldopt.mapFn(fv.Addr(), vSrc); err != nil {
					return fmt.Errorf("mapping field %q: %w", fldopt.dst.Name, err)
				}
			}

			return nil
		},
	}
}

type AutoTypeOptions struct {
	dst reflect.Type
	src reflect.Type

	fldOpts []*AutoFieldOptions
}

func (o *AutoTypeOptions) Dst() reflect.Type {
	return o.dst
}

func (o *AutoTypeOptions) Src() reflect.Type {
	return o.src
}

func (o *AutoTypeOptions) Field(name string, fn interface{}) {
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
		panic(fmt.Errorf("fn function must have 1 return value, but had %d", t.NumOut()))
	}

	var opts AutoFieldOptions
	switch t.NumIn() {
	case 2:
		opts.name = name

		sf, found := o.dst.Elem().FieldByName(name)
		if !found {
			panic(fmt.Errorf("field %q does not exist on %q", name, o.dst))
		}

		opts.dst = sf

		v := reflect.ValueOf(fn)
		opts.mapFn = func(dst reflect.Value, src reflect.Value) error {
			result := v.Call([]reflect.Value{
				dst,
				src,
			})

			if result[0].IsNil() {
				return nil
			}

			return result[0].Interface().(error)
		}

		o.fldOpts = append(o.fldOpts, &opts)
	default:
		panic(fmt.Errorf("fn function must have 2 or 3 arguments, but had %d", t.NumIn()))
	}
}

type AutoFieldOptions struct {
	name string
	dst reflect.StructField

	mapFn func(reflect.Value, reflect.Value) error
}
