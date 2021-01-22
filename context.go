package mapper

type Context interface {
	Map(dst interface{}, src interface{}) error
}
