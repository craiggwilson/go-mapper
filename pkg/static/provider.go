package static

import (
	"github.com/craiggwilson/go-mapper/pkg/core"
)

// NewProvider makes a Provider.
func NewProvider() *Provider {
	return &Provider{}
}

// Provider holds a static list of type mappers.
type Provider struct {
	tms []core.Mapper
}

// Mappers implements the core.Provider interface.
func (p *Provider) Mappers() []core.Mapper {
	return p.tms
}

// Add adds a Mapper to the static mapping list.
func (p *Provider) Add(tm core.Mapper) {
	p.tms = append(p.tms, tm)
}
