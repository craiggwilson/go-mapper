package naming

import (
	"strings"
	"unicode"
)

// Pascal returns the implementation for a PascalCase naming strategy.
func Pascal() PascalCase {
	return PascalCase{}
}

// PascalCaseNamingConvention returns possible matches based on PascalCase splitting.
type PascalCase struct{}

// Possibilities returns the possible matches for the given name.
func (PascalCase) Possibilities(name string) []Possibility {
	parts := splitFunc(name, unicode.IsUpper)
	var results []Possibility
	for i := 0; i < len(parts); i++ {
		results = append(results, Possibility{
			Match: strings.Join(parts[:len(parts)-i], ""),
			Remaining: strings.Join(parts[len(parts)-i:], "")})
	}
	return results
}
