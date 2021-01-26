package naming

// Exact returns the implementation for an exact match naming strategy.
func Exact() ExactMatch {
	return ExactMatch{}
}

// ExactMatch provides only an exact match possibility.
type ExactMatch struct {}

// Possibilities returns the match possibilities for the given name.
func (ExactMatch) Possibilities(name string) []Possibility {
	return []Possibility{{name, ""}}
}
