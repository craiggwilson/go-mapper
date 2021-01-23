package auto

import (
	"fmt"
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/core"
)

// NewProvider makes an Provider.
func NewProvider() *Provider {
	return &Provider{}
}

// Provider is used to automatically map types following prescribed strategies for naming and type conversion.
type Provider struct {
	// strategies

	mappers []core.Mapper
}

// Mappers implements the core.Provider interface.
func (c *Provider) Mappers() []core.Mapper {
	return c.mappers
}

// AddStruct registers a struct for mapping. The fn argument must match the signature
// func(dst <type>, src <type>) or func(dst <type>, src <type>, cfg *StructOptions). If fn is not a function,
// or it's signature does not match the requirements, a panic is raised.
func (c *Provider) AddStruct(fn interface{}) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("fn argument must be a func but got a %q", t.Kind()))
	}

	switch t.NumOut() {
	case 0:
	default:
		panic(fmt.Errorf("fn function must have no return values, but had %d", t.NumOut()))
	}

	opts := newStructOptions()
	switch t.NumIn() {
	case 3:
		if !t.In(2).AssignableTo(tAutoTypeConfig) {
			panic(fmt.Errorf("fn function with 3 arguments must have *StructOptions as the last, but got %q", t.In(2)))
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
			reflect.ValueOf(opts),
		})

	case 2:
		opts.dst = t.In(0)
		opts.src = t.In(1)
	default:
		panic(fmt.Errorf("fn function must have 2 or 3 arguments, but had %d", t.NumIn()))
	}

	dstStruct := opts.dst.Elem()

	for i := 0; i < dstStruct.NumField(); i++ {
		fld := dstStruct.Field(i)
		if _, ok := opts.fields[fld.Name]; ok {
			continue
		}

		mapFn := opts.fieldMappingStrategy.Create(fld, opts.src)
		if mapFn != nil {
			opts.fields[fld.Name] = &FieldOptions{
				dst: fld,
				mapFn: mapFn,
			}
		}
	}

	tm := c.createMapper(opts)
	c.mappers = append(c.mappers, tm)
}

func (c *Provider) createMapper(opts *StructOptions) core.Mapper {
	return core.NewFunctionMapper(opts.dst, opts.src, func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
		dst = reflect.Indirect(dst)
		for _, fld := range opts.fields {
			fv := dst.FieldByIndex(fld.dst.Index)
			if !fv.CanAddr() {
				return fmt.Errorf("field %q cannot be addressed", fld.dst.Name)
			}
			if err := fld.mapFn(ctx, fv.Addr(), src); err != nil {
				return fmt.Errorf("mapping field %q: %w", fld.dst.Name, err)
			}
		}

		return nil
	})
}

func newStructOptions() *StructOptions {
	return &StructOptions{
		fieldMappingStrategy: ExactNameFieldMappingStrategy{},
		fields: make(map[string]*FieldOptions),
	}
}

type StructOptions struct {
	dst reflect.Type
	src reflect.Type

	fieldMappingStrategy FieldMappingStrategy

	fields map[string]*FieldOptions
}

func (o *StructOptions) Dst() reflect.Type {
	return o.dst
}

func (o *StructOptions) Src() reflect.Type {
	return o.src
}

func (o *StructOptions) Field(name string, fn interface{}) {
	fm := core.MapperFromFunc(fn)
	sf, found := o.dst.Elem().FieldByName(name)
	if !found {
		panic(fmt.Errorf("field %q does not exist on %q", name, o.dst))
	}

	if !reflect.PtrTo(sf.Type).AssignableTo(fm.Dst()) {
		panic(fmt.Errorf("dst argument must be assignable from field %q: have %q but need %q", name, fm.Dst(), sf.Type))
	}

	var opts FieldOptions
	opts.mapFn = fm.Func()
	opts.dst = sf

	o.fields[sf.Name] = &opts
}

// FieldOptions contains options for a field mapping.
type FieldOptions struct {
	dst reflect.StructField

	mapFn core.MapperFunc
}

// FieldMappingStrategy is a strategy for automatically mapping fields.
type FieldMappingStrategy interface {
	// Create makes a core.MapperFunc for the given destination. If one cannot be made, then nil should be returned.
	Create(dst reflect.StructField, src reflect.Type) core.MapperFunc
}

type ExactNameFieldMappingStrategy struct {}

func (ExactNameFieldMappingStrategy) Create(dst reflect.StructField, src reflect.Type) core.MapperFunc {
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	if src.Kind() != reflect.Struct {
		return nil
	}

	srcField, found := src.FieldByName(dst.Name)
	if !found {
		return nil
	}

	// TODO: ensure types are compatible

	return func(ctx core.Context, vDst reflect.Value, vSrc reflect.Value) error {
		if vSrc.IsNil() {
			return nil
		}
		vSrc = reflect.Indirect(vSrc).FieldByIndex(srcField.Index)
		vDst.Elem().Set(vSrc)
		return nil
	}
}
