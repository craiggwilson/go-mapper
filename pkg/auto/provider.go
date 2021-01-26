package auto

import (
	"fmt"
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/core"
	"github.com/craiggwilson/go-mapper/pkg/reflecth"
)

// NewProvider makes an Provider.
func NewProvider() *Provider {
	return &Provider{
		converterFactory: reflecth.ConverterFactoryFunc(reflecth.ConverterFor),
		namingConvention: PascalCaseNamingConvention{},
	}
}

// Provider is used to automatically map types following prescribed strategies for naming and type conversion.
type Provider struct {
	// strategies
	converterFactory reflecth.ConverterFactory
	namingConvention NamingConvention

	opts []*StructOptions
}

// Mappers implements the core.Provider interface.
func (p *Provider) Mappers() []core.Mapper {
	mappers := make([]core.Mapper, 0, len(p.opts))
	for _, opt := range p.opts {
		mappers = append(mappers, p.createMapper(opt))
	}

	return mappers
}

// WithConverterFactory applies the converterFactory to all future uses.
func (p *Provider) WithConverterFactory(cf reflecth.ConverterFactory) {
	if cf == nil {
		panic(fmt.Errorf("cf cannot be nil"))
	}

	p.converterFactory = cf
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
		fields:           make(map[string]*FieldOptions),
		converterFactory: p.converterFactory,
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

	p.opts = append(p.opts, opts)
}

func (p *Provider) createMapper(opts *StructOptions) core.Mapper {
	// set defaults if necessary
	if opts.converterFactory == nil {
		opts.converterFactory = p.converterFactory
	}
	if opts.namingConvention == nil {
		opts.namingConvention = p.namingConvention
	}

	// AutoMap any remaining fields
	dstStruct := opts.dst.Elem()
	for i := 0; i < dstStruct.NumField(); i++ {
		fld := dstStruct.Field(i)
		if _, ok := opts.fields[fld.Name]; ok {
			continue
		}

		accessor := matchNameToSource(opts.namingConvention, fld.Name, opts.src)
		if accessor == nil {
			continue
		}

		converter, err := opts.converterFactory.ConverterFor(fld.Type, accessor.Type())
		if err != nil {
			continue
		}

		opts.fields[fld.Name] = &FieldOptions{
			dst: fld,
			accessor: accessor,
			converter: converter,
			mapFn: func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
				if src.IsNil() {
					return nil
				}

				src = accessor.ValueFrom(src)

				if converter != nil {
					return converter.Convert(dst, src)
				}

				dst.Elem().Set(src)
				return nil
			},
		}
	}

	return core.NewFunctionMapper(opts.dst, opts.src, func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
		dst = reflect.Indirect(dst)
		for _, fld := range opts.fields {
			fv := dst.FieldByIndex(fld.dst.Index)
			if !fv.CanAddr() {
				return fmt.Errorf("field %q cannot be addressed", fld.dst.Name)
			}

			if err := fld.mapFn(ctx, fv.Addr(), src); err != nil {
				name := "(custom function)"
				if fld.accessor != nil {
					name = fld.accessor.Name()
				}
				return fmt.Errorf("mapping field %q from %q: %w",
					fmt.Sprintf("%v.%s", dst.Type(), fld.dst.Name),
					fmt.Sprintf("%v.%s", src.Type(), name),
					err)
			}
		}

		return nil
	})
}

type StructOptions struct {
	dst reflect.Type
	src reflect.Type

	converterFactory reflecth.ConverterFactory
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

func (o *StructOptions) WithConverter(cf reflecth.ConverterFactory) {
	if cf == nil {
		panic(fmt.Errorf("cf cannot be nil"))
	}

	o.converterFactory = cf
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

	converter reflecth.Converter
	accessor reflecth.Accessor

	mapFn core.MapperFunc
}
