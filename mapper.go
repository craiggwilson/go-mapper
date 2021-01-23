package mapper

import (
	"reflect"
)

// New makes a Mapper.
func New(providers ...TypeMapperProvider) (*Mapper, error) {
	src := make(map[reflect.Type]map[reflect.Type]TypeMapper)
	for _, p := range providers {
		for _, tm := range p.TypeMappers() {
			dstMap, ok := src[tm.Src()]
			if !ok {
				dstMap = make(map[reflect.Type]TypeMapper)
				src[tm.Src()] = dstMap
			}

			_, ok = dstMap[tm.Dst()]
			if ok {
				return nil, newDuplicateTypeMapError(tm.Dst(), tm.Src())
			}

			dstMap[tm.Dst()] = tm
		}
	}

	return &Mapper{
		src: src,
	}, nil
}

type Mapper struct {
	src map[reflect.Type]map[reflect.Type]TypeMapper
}

func (m *Mapper) Map(dst interface{}, src interface{}) error {
	tDst := reflect.TypeOf(dst)
	tSrc := reflect.TypeOf(src)
	if dstMap, ok := m.src[tSrc]; ok {
		if tm, ok := dstMap[tDst]; ok {
			return tm.Map(nil, dst, src)
		}
	}

	return ErrNoTypeMapperFound
}

// TypeMapperProvider provides TypeMappers.
type TypeMapperProvider interface {
	TypeMappers() []TypeMapper
}

// TypeMapper handles mapping from src to dst.
type TypeMapper interface {
	// Dst is the type of the destination.
	Dst() reflect.Type
	// Src is the type of the source.
	Src() reflect.Type
	// Map performs the mapping to dst from src.
	Map(ctx Context, dst interface{}, src interface{}) error
}
