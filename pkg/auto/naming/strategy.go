package naming

// Strategy represents a strategy for discovering possible matches based on a given name.
type Strategy interface {
	// Next returns the next possible matches given the name.
	Possibilities(name string) []Possibility
}
