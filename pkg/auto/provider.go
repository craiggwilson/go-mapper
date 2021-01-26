package auto

import (
	"fmt"
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/auto/naming"
	"github.com/craiggwilson/go-mapper/pkg/core"
	"github.com/craiggwilson/go-mapper/pkg/reflecth"
)

// NewProvider makes an Provider.
func NewProvider() *Provider {
	return &Provider{
		converterFactory: reflecth.ConverterFactoryFunc(reflecth.ConverterFor),
		namingStrategy:   naming.PascalCase{},
	}
}

// Provider is used to automatically map types following prescribed strategies for naming and type conversion.
type Provider struct {
	// strategies
	converterFactory reflecth.ConverterFactory
	namingStrategy   naming.Strategy

	structs []*Struct
}

// Mappers implements the core.Provider interface.
func (p *Provider) Mappers() []core.Mapper {
	mappers := make([]core.Mapper, 0, len(p.structs))
	for _, opt := range p.structs {
		mappers = append(mappers, p.createMapper(opt))
	}

	return mappers
}

// WithConverterFactory applies the converterFactory to all future uses.
func (p *Provider) WithConverterFactory(cf reflecth.ConverterFactory) {
	p.converterFactory = cf
}

// WithNamingConvention applies the naming convention to all future uses.
func (p *Provider) WithNamingConvention(ns naming.Strategy) {
	p.namingStrategy = ns
}

// Add adds a src and dst to automatically create a core.Mapper.
func (p *Provider) Add(dst reflect.Type, src reflect.Type, opts ...func(structOpts)) {
	s := Struct{
		dst: dst,
		src: src,
		fields: make(map[string]*Field),
	}

	for _, opt := range opts {
		opt(&s)
	}

	p.structs = append(p.structs, &s)
}

func (p *Provider) createMapper(s *Struct) core.Mapper {
	converterFactory := s.converterFactory
	if converterFactory == nil {
		converterFactory = p.converterFactory
	}
	namingStrategy := s.namingStrategy
	if namingStrategy == nil {
		namingStrategy = p.namingStrategy
	}

	dst := s.dst.Elem()
	fields := make(map[string]*Field, dst.NumField())
	for k, fld := range s.fields {
		fields[k] = fld
	}

	for i := 0; i < dst.NumField(); i++ {
		fld := dst.Field(i)
		f, ok := fields[fld.Name]
		if ok && f.mapper != nil {
			// If we have a mapper already, we don't need to do any automapping work.
			continue
		} else if !ok {
			f = &Field{
				dst: fld,
			}
		}

		accessor := f.accessor
		if accessor == nil {
			accessor = matchNameToSource(namingStrategy, fld.Name, s.src)
			if accessor == nil {
				continue
			}
		}

		converter := f.converter
		if converter == nil {
			var err error
			converter, err = converterFactory.ConverterFor(fld.Type, accessor.Type())
			if err != nil {
				continue
			}
		}

		fields[fld.Name] = &Field{
			dst: fld,
			accessor: accessor,
			converter: converter,
			mapper: core.NewFunctionMapper(
				fld.Type,
				s.src,
				func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
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
			),
		}
	}

	return core.NewFunctionMapper(
		s.dst,
		s.src,
		func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
			dst = reflect.Indirect(dst)
			for _, fld := range fields {
				fv := dst.FieldByIndex(fld.dst.Index)
				if !fv.CanAddr() {
					return fmt.Errorf("field %q cannot be addressed", fld.dst.Name)
				}

				if err := fld.mapper.Map(ctx, fv.Addr(), src); err != nil {
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
		},
	)
}

type Struct struct {
	dst reflect.Type
	src reflect.Type

	converterFactory reflecth.ConverterFactory
	namingStrategy   naming.Strategy

	fields map[string]*Field
}

func (s *Struct) Dst() reflect.Type {
	return s.dst
}

func (s *Struct) Src() reflect.Type {
	return s.src
}

func (s *Struct) WithConverterFactory(cf reflecth.ConverterFactory) {
	s.converterFactory = cf
}

func (s *Struct) WithField(name string, opts ...func(fieldOpts)) {
	sf, found := s.dst.Elem().FieldByName(name)
	if !found {
		panic(fmt.Errorf("field %q does not exist on %v", name, s.dst))
	}

	f := Field{
		dst: sf,
	}
	for _, opt := range opts {
		opt(&f)
	}

	s.fields[sf.Name] = &f
}

func (s *Struct) WithNamingStrategy(ns naming.Strategy) {
	s.namingStrategy = ns
}

// Field contains options for a field mapping.
type Field struct {
	dst reflect.StructField

	accessor       reflecth.Accessor
	converter      reflecth.Converter
	mapper         core.Mapper
	namingStrategy naming.Strategy
}

func (f *Field) WithAccessor(a reflecth.Accessor) {
	f.accessor = a
}

func (f *Field) WithConverter(c reflecth.Converter) {
	f.converter = c
}

func (f *Field) WithMapper(m core.Mapper) {
	f.mapper = m
}

func (f *Field) WithNamingStrategy(ns naming.Strategy) {
	f.namingStrategy = ns
}
