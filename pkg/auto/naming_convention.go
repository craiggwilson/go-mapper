package auto

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/craiggwilson/go-mapper/pkg/reflecth"
)

func matchNameToSource(nc NamingConvention, name string, src reflect.Type) reflecth.Accessor {
	var currentAccessor reflecth.Accessor
	lastName := name
	currentName := name
	currentType := src
	for len(currentName) > 0 {
		for currentType.Kind() == reflect.Ptr {
			currentType = currentType.Elem()
		}

		for _, pm := range nc.Next(currentName) {
			fld, found := currentType.FieldByName(pm[0])
			if !found {
				continue
			}

			if currentAccessor == nil {
				currentAccessor = reflecth.NewFieldAccessor(fld)
			} else {
				currentAccessor = reflecth.NewAccessorPair(currentAccessor, reflecth.NewFieldAccessor(fld))
			}

			lastName = currentName
			currentName = pm[1]
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

type NamingConvention interface {
	Next(name string) [][2]string
}

type ExactMatchNamingConvention struct {}

func (ExactMatchNamingConvention) Next(name string) [][2]string {
	return [][2]string{{name, ""}}
}

// PascalCaseNamingConvention returns possible matches based on PascalCase splitting.
type PascalCaseNamingConvention struct{}

func (PascalCaseNamingConvention) Next(name string) [][2]string {
	parts := splitFunc(name, unicode.IsUpper)
	var results [][2]string
	for i := 0; i < len(parts); i++ {
		results = append(results, [2]string{strings.Join(parts[:len(parts)-i], ""), strings.Join(parts[len(parts)-i:], "")})
	}
	return results
}

// adapted version of strings.FieldsFunc.
func splitFunc(s string, f func(rune) bool) []string {
	type span struct {
		start int
		end   int
	}
	spans := make([]span, 0, 4)

	start := 0
	for end, r := range s {
		if f(r) && start < end {
			spans = append(spans, span{start, end})
			start = end
		}
	}

	if start < len(s) {
		spans = append(spans, span{start, len(s)})
	}

	a := make([]string, len(spans))
	for i, span := range spans {
		a[i] = s[span.start:span.end]
	}
	return a
}