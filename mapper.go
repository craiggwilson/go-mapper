package mapper

import (
	"reflect"
)


type Mapper struct {
	providers []TypeMapperProvider
}

func (m *Mapper) Map(dst interface{}, src interface{}) error {
	for _, p := range m.providers {
		tDst := reflect.TypeOf(dst)
		tSrc := reflect.TypeOf(src)
		tm, err := p.TypeMapperFor(tDst, tSrc)
		if err == ErrNoTypeMapperFound {
			continue
		}

		if err != nil {
			return err
		}

		return tm.Map(nil, dst, src)
	}

	return ErrNoTypeMapperFound
}
