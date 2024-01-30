package utility

type Option[T any] interface {
	apply(*T)
}

type OptionFunc[T any] func(*T)

func (f OptionFunc[T]) apply(t *T) {
	f(t)
}

func Apply[T any](o *T, opts ...Option[T]) {
	for _, opt := range opts {
		opt.apply(o)
	}
}
