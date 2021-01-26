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
func (p *Provider) WithNamingConvention(nc NamingConvention) {
	p.namingConvention = nc
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
	namingConvention := s.namingConvention
	if namingConvention == nil {
		namingConvention = p.namingConvention
	}

	dstStruct := s.dst.Elem()
	for i := 0; i < dstStruct.NumField(); i++ {
		fld := dstStruct.Field(i)
		f, ok := s.fields[fld.Name]
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
			accessor = matchNameToSource(namingConvention, fld.Name, s.src)
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

		s.fields[fld.Name] = &Field{
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

	return core.NewFunctionMapper(s.dst, s.src, func(ctx core.Context, dst reflect.Value, src reflect.Value) error {
		dst = reflect.Indirect(dst)
		for _, fld := range s.fields {
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
	})
}

type Struct struct {
	dst reflect.Type
	src reflect.Type

	converterFactory reflecth.ConverterFactory
	namingConvention NamingConvention

	fields map[string]*Field
}

func (s *Struct) Dst() reflect.Type {
	return s.dst
}

func (s *Struct) Src() reflect.Type {
	return s.src
}

func (s *Struct) Field(name string, fn interface{}) {
	fm := core.MapperFromFunc(fn)
	sf, found := s.dst.Elem().FieldByName(name)
	if !found {
		panic(fmt.Errorf("field %q does not exist on %q", name, s.dst))
	}

	if !reflect.PtrTo(sf.Type).AssignableTo(fm.Dst()) {
		panic(fmt.Errorf("dst argument must be assignable from field %q: have %q but need %q", name, fm.Dst(), sf.Type))
	}

	var opts Field
	opts.mapper = fm
	opts.dst = sf

	s.fields[sf.Name] = &opts
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

func (s *Struct) WithNamingConvention(nc NamingConvention) {
	s.namingConvention = nc
}

// Field contains options for a field mapping.
type Field struct {
	dst reflect.StructField


	accessor reflecth.Accessor
	converter reflecth.Converter
	mapper core.Mapper
	namingConvention NamingConvention
}

func (f *Field) WithConverter(c reflecth.Converter) {
	f.converter = c
}

func (f *Field) WithMapper(m core.Mapper) {
	f.mapper = m
}

func (f *Field) WithNamingConvention(nc NamingConvention) {
	f.namingConvention = nc
}


//// AddStruct registers a struct for mapping. The fn argument must match the signature
//// func(dst <type>, src <type>) or func(dst <type>, src <type>, cfg *Struct). If fn is not a function,
//// or it's signature does not match the requirements, a panic is raised.
//func (p *Provider) AddStruct(fn interface{}) {
//	t := reflect.TypeOf(fn)
//	if t.Kind() != reflect.Func {
//		panic(fmt.Sprintf("fn argument must be a func but got a %q", t.Kind()))
//	}
//
//	switch t.NumOut() {
//	case 0:
//	default:
//		panic(fmt.Errorf("fn function must have no return values, but had %d", t.NumOut()))
//	}
//
//	opts := &Struct{
//		fields:           make(map[string]*Field),
//		converterFactory: p.converterFactory,
//		namingConvention: p.namingConvention,
//	}
//
//	switch t.NumIn() {
//	case 3:
//		if !t.In(2).AssignableTo(tAutoTypeConfig) {
//			panic(fmt.Errorf("fn function with 3 arguments must have *Struct as the last, but got %q", t.In(2)))
//		}
//
//		opts.dst = t.In(0)
//		if opts.dst.Kind() != reflect.Ptr || opts.dst.Elem().Kind() != reflect.Struct {
//			panic(fmt.Errorf("fn function's first argument must be a pointer to a struct"))
//		}
//		opts.src = t.In(1)
//
//		v := reflect.ValueOf(fn)
//
//		_ = v.Call([]reflect.Value{
//			reflect.Zero(opts.dst),
//			reflect.Zero(opts.src),
//			reflect.ValueOf(opts),
//		})
//
//	case 2:
//		opts.dst = t.In(0)
//		opts.src = t.In(1)
//	default:
//		panic(fmt.Errorf("fn function must have 2 or 3 arguments, but had %d", t.NumIn()))
//	}
//
//	p.structs = append(p.structs, opts)
//}