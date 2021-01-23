package mapper

// NewStaticConfig makes a StaticConfig.
func NewStaticConfig() *StaticConfig {
	return &StaticConfig{}
}

// StaticConfig holds a static list of type mappers.
type StaticConfig struct {
	tms []TypeMapper
}

// TypeMapperFor implements the TypeMapperProvider interface.
func (p *StaticConfig) TypeMappers() []TypeMapper{
	return p.tms
}

// Add adds a TypeMapper to the static mapping list.
func (p *StaticConfig) Add(tm TypeMapper) {
	p.tms = append(p.tms, tm)
}
