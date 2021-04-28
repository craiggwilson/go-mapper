package core

// Provider provides Mappers.
type Provider interface {
	Mappers() ([]Mapper, error)
}
