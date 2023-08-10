package opts

type Option[T any] interface {
	Apply(cfg *T)
}

type OptionFunc[T any] func(cfg *T)

func (f OptionFunc[T]) Apply(cfg *T) {
	f(cfg)
}
