package auto

import (
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/auto/naming"
	"github.com/craiggwilson/go-mapper/pkg/reflecth"
)

func matchNameToSource(nc naming.Strategy, name string, src reflect.Type) reflecth.Accessor {
	var currentAccessor reflecth.Accessor
	lastName := name
	currentName := name
	currentType := src
	for len(currentName) > 0 {
		for currentType.Kind() == reflect.Ptr {
			currentType = currentType.Elem()
		}

		for _, p := range nc.Possibilities(currentName) {
			fld, found := currentType.FieldByName(p.Match)
			if !found {
				continue
			}

			if currentAccessor == nil {
				currentAccessor = reflecth.NewFieldAccessor(fld)
			} else {
				currentAccessor = reflecth.NewAccessorPair(currentAccessor, reflecth.NewFieldAccessor(fld))
			}

			lastName = currentName
			currentName = p.Remaining
			currentType = currentAccessor.Type()
			break
		}

		// we made no progress, so we are done
		if lastName == currentName {
			return nil
		}
	}

	return currentAccessor
}
