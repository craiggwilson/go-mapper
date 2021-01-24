package auto

import (
	"fmt"
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/auto/convert"
	"github.com/craiggwilson/go-mapper/pkg/core"
)

// NewProvider makes an Provider.
func NewProvider() *Provider {
	return &Provider{
		namingConvention: PascalCaseNamingConvention{},
	}
}

// Provider is used to automatically map types following prescribed strategies for naming and type conversion.
type Provider struct {
	// strategies
	namingConvention NamingConvention
	converter        convert.Converter

	opts []*StructOptions

	mappers []core.Mapper
}

// Mappers implements the core.Provider interface.
func (p *Provider) Mappers() []core.Mapper {
	return p.mappers
}

// WithConverter applies the converter to all future uses.
func (p *Provider) WithConverter(c convert.Converter) {
	if c == nil {
		panic(fmt.Errorf("c cannot be nil"))
	}

	p.converter = c
}

// WithNamingConvention applies the naming convention to all future uses.
func (p *Provider) WithNamingConvention(nc NamingConvention) {
	if nc == nil {
		panic(fmt.Errorf("nc cannot be nil"))
	}

	p.namingConvention = nc
}

// AddStruct registers a struct for mapping. The fn argument must match the signature
// func(dst <type>, src <type>) or func(dst <type>, src <type>, cfg *StructOptions). If fn is not a function,
// or it's signature does not match the requirements, a panic is raised.
func (p *Provider) AddStruct(fn interface{}) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("fn argument must be a func but got a %q", t.Kind()))
	}

	switch t.NumOut() {
	case 0:
	default:
		panic(fmt.Errorf("fn function must have no return values, but had %d", t.NumOut()))
	}

	opts := &StructOptions{
		fields: make(map[string]*FieldOptions),
		converter: p.converter,
		namingConvention: p.namingConvention,
	}

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

		accessor := applyNamingConvention(opts.namingConvention, fld.Name, opts.src)
		if accessor == nil {
			continue
		}

		converter := opts.converter

		opts.fields[fld.Name] = &FieldOptions{
			dst: fld,
			mapFn: func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
				if src.IsNil() {
					return nil
				}
				src = accessor.ValueFrom(src)

				dst = dst.Elem()
				if !src.Type().AssignableTo(dst.Type()) {
					return converter.Convert(dst, src)
				}

				dst.Set(src)
				return nil
			},
		}
	}

	tm := p.createMapper(opts)
	p.mappers = append(p.mappers, tm)
}

func (p *Provider) createMapper(opts *StructOptions) core.Mapper {
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

type StructOptions struct {
	dst reflect.Type
	src reflect.Type

	converter convert.Converter
	namingConvention NamingConvention

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

func (o *StructOptions) WithConverter(c convert.Converter) {
	if c == nil {
		panic(fmt.Errorf("c cannot be nil"))
	}

	o.converter = c
}

func (o *StructOptions) WithNamingConvention(nc NamingConvention) {
	if nc == nil {
		panic(fmt.Errorf("nc cannot be nil"))
	}

	o.namingConvention = nc
}

// FieldOptions contains options for a field mapping.
type FieldOptions struct {
	dst reflect.StructField

	mapFn core.MapperFunc
}
