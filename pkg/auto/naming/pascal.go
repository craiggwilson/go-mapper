package naming

import (
	"strings"
	"unicode"
)

// PascalCaseNamingConvention returns possible matches based on PascalCase splitting.
type PascalCase struct{}

// Possibilities returns the possibile matches for the given name.
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
