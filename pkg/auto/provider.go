package auto

import (
	"fmt"
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/auto/accessor"
	"github.com/craiggwilson/go-mapper/pkg/auto/converter"
	"github.com/craiggwilson/go-mapper/pkg/auto/naming"
	"github.com/craiggwilson/go-mapper/pkg/core"
	"github.com/craiggwilson/go-mapper/pkg/internal"
)

// NewProvider makes an Provider.
func NewProvider() *Provider {
	return &Provider{
		converterFactory: converter.FactoryFunc(converter.For),
		namingStrategy:   naming.PascalCase{},
	}
}

// Provider is used to automatically map types following prescribed strategies for naming and type conversion.
type Provider struct {
	// strategies
	converterFactory converter.Factory
	namingStrategy   naming.Strategy

	structs []*Struct
}

// Mappers implements the core.Provider interface.
func (p *Provider) Mappers() ([]core.Mapper, error) {
	mappers := make([]core.Mapper, 0, len(p.structs))
	for _, opt := range p.structs {
		mapper, err := p.createMapper(opt)
		if err != nil {
			return nil, err
		}
		mappers = append(mappers, mapper)
	}

	return mappers, nil
}

// WithConverterFactory applies the converterFactory to all future uses.
func (p *Provider) WithConverterFactory(cf converter.Factory) {
	p.converterFactory = cf
}

// WithNamingConvention applies the naming convention to all future uses.
func (p *Provider) WithNamingConvention(ns naming.Strategy) {
	p.namingStrategy = ns
}

// Add adds a src and dst to automatically create a core.Mapper.
func (p *Provider) Add(dst reflect.Type, src reflect.Type, opts ...func(structOpts)) {
	// Use the bare, non-pointer types for mapping.
	dst = internal.UnwrapPtrType(dst)
	src = internal.UnwrapPtrType(src)

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

func (p *Provider) createMapper(s *Struct) (core.Mapper, error) {
	converterFactory := s.converterFactory
	if converterFactory == nil {
		converterFactory = p.converterFactory
	}
	namingStrategy := s.namingStrategy
	if namingStrategy == nil {
		namingStrategy = p.namingStrategy
	}

	fields := make(map[string]*Field, s.dst.NumField())
	for k, fld := range s.fields {
		fields[k] = fld
	}

	for i := 0; i < s.dst.NumField(); i++ {
		fld := s.dst.Field(i)
		f, ok := fields[fld.Name]
		if ok && f.mapper != nil {
			// If we have a mapper already, we don't need to do any automapping work.
			continue
		} else if !ok {
			f = &Field{
				dst: fld,
			}
		}

		if f.ignore {
			continue
		}

		acc := f.accessor
		if acc == nil {
			acc = findAccessor(namingStrategy, fld.Name, s.src)
			if acc == nil {
				continue
			}
		}

		conv := f.converter
		if conv == nil {
			var err error
			conv, err = converterFactory.ConverterFor(fld.Type, acc.Type())
			if err != nil {
				return nil, fmt.Errorf("mapping field %q from %q: %w",
					fmt.Sprintf("%v.%s", s.dst.Name(), fld.Name),
					fmt.Sprintf("%v.%s", s.src.Name(), acc.Name()),
					err)
			}
		}

		fields[fld.Name] = &Field{
			dst:       fld,
			accessor:  acc,
			converter: conv,
			mapper: core.NewFunctionMapper(
				fld.Type,
				s.src,
				func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
					if src.IsNil() {
						return nil
					}

					src = acc.ValueFrom(src)

					if conv != nil {
						return conv.Convert(dst, src)
					}

					dst = internal.EnsureSettableDst(dst)
					dst.Set(internal.UnwrapPtrValue(src))
					return nil
				},
			),
		}
	}

	return core.NewFunctionMapper(
		s.dst,
		s.src,
		func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
			dst = internal.EnsureSettableDst(dst)
			for _, fld := range fields {
				fv := dst.FieldByIndex(fld.dst.Index)
				if fld.ignore {
					continue
				}

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
	), nil
}

type Struct struct {
	dst reflect.Type
	src reflect.Type

	converterFactory converter.Factory
	namingStrategy   naming.Strategy

	fields map[string]*Field
}

func (s *Struct) Dst() reflect.Type {
	return s.dst
}

func (s *Struct) Src() reflect.Type {
	return s.src
}

func (s *Struct) WithConverterFactory(cf converter.Factory) {
	s.converterFactory = cf
}

func (s *Struct) WithField(name string, opts ...func(fieldOpts)) {
	sf, found := s.dst.FieldByName(name)
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

	accessor       accessor.Accessor
	converter      converter.Converter
	ignore bool
	mapper         core.Mapper
	namingStrategy naming.Strategy
}

func (f *Field) WithAccessor(a accessor.Accessor) {
	f.accessor = a
}

func (f *Field) WithConverter(c converter.Converter) {
	f.converter = c
}

func (f *Field) WithIgnore() {
	f.ignore = true
}

func (f *Field) WithMapper(m core.Mapper) {
	f.mapper = m
}

func (f *Field) WithNamingStrategy(ns naming.Strategy) {
	f.namingStrategy = ns
}
