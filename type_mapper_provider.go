package mapper

import (
	"reflect"
)

// TypeMapperProvider return TypeMappers for the given combination of src to dst types.
type TypeMapperProvider interface {
	// TypeMapperFor returns a TypeMapper for mapping between dst and src. If the provider cannot provide one,
	// ErrNoTypeMapperFound should be returned.
	TypeMapperFor(dst reflect.Type, src reflect.Type) (TypeMapper, error)
}

type staticTypeMapperProvider struct {
	src map[reflect.Type]map[reflect.Type]TypeMapper
}

func (p *staticTypeMapperProvider) TypeMapperFor(dst reflect.Type, src reflect.Type) (TypeMapper, error) {
	if dstMap, ok := p.src[src]; ok {
		if tm, ok := dstMap[dst]; ok {
			return tm, nil
		}
	}

	return nil, ErrNoTypeMapperFound
}

func (p *staticTypeMapperProvider) addTypeMapper(tm TypeMapper) error {
	if p.src == nil {
		p.src = make(map[reflect.Type]map[reflect.Type]TypeMapper)
	}

	dstMap, ok := p.src[tm.Src()]
	if !ok {
		dstMap = make(map[reflect.Type]TypeMapper)
		p.src[tm.Src()] = dstMap
	}

	_, ok = dstMap[tm.Dst()]
	if ok {
		return newDuplicateTypeMapError(tm.Dst(), tm.Src())
	}

	dstMap[tm.Dst()] = tm
	return nil
}




