package naming

// ExactMatch provides only an exact match possibility.
type ExactMatch struct {}

// Possibilities returns the match possibilities for the given name.
func (ExactMatch) Possibilities(name string) []Possibility {
	return []Possibility{{name, ""}}
}
