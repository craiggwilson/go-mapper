package auto

import (
	"reflect"

	"github.com/craiggwilson/go-mapper/pkg/auto/accessor"
	"github.com/craiggwilson/go-mapper/pkg/auto/naming"
)

func findAccessor(ns naming.Strategy, name string, src reflect.Type) accessor.Accessor {
	var currentAccessor accessor.Accessor
	lastName := name
	currentName := name
	currentType := src
	for len(currentName) > 0 {
		for currentType.Kind() == reflect.Ptr {
			currentType = currentType.Elem()
		}

		for _, p := range ns.Possibilities(currentName) {
			fld, found := currentType.FieldByName(p.Match)
			if !found {
				continue
			}

			if currentAccessor == nil {
				currentAccessor = accessor.Field(fld)
			} else {
				currentAccessor = accessor.Pair(currentAccessor, accessor.Field(fld))
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
