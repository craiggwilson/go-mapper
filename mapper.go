package mapper

import (
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/core"
)

// New makes a Mapper.
func New(providers ...core.Provider) (*Mapper, error) {
	src := make(map[reflect.Type]map[reflect.Type]core.Mapper)
	for _, p := range providers {
		mappers, err := p.Mappers()
		if err != nil {
			return nil, err
		}
		for _, tm := range mappers {
			dstMap, ok := src[tm.Src()]
			if !ok {
				dstMap = make(map[reflect.Type]core.Mapper)
				src[tm.Src()] = dstMap
			}

			dstMap[tm.Dst()] = tm
		}
	}

	return &Mapper{
		src: src,
	}, nil
}

type Mapper struct {
	src map[reflect.Type]map[reflect.Type]core.Mapper
}

func (m *Mapper) Map(dst interface{}, src interface{}) error {
	vDst := reflect.ValueOf(dst)
	vSrc := reflect.ValueOf(src)
	if dstMap, ok := m.src[vSrc.Type()]; ok {
		if tm, ok := dstMap[vDst.Type()]; ok {
			return tm.Map(nil, vDst, vSrc)
		}
	}

	return ErrNoTypeMapperFound
}


