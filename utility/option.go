package utility

type Option[T any] interface {
	apply(*T) error
}

type OptionFunc[T any] func(*T) error

func (f OptionFunc[T]) apply(t *T) error {
	return f(t)
}

func InlineOpt[T any](f func(*T)) Option[T] {
	return OptionFunc[T](func(t *T) error {
		f(t)
		return nil
	})
}

func Apply[T any](o *T, opts ...Option[T]) error {
	for _, opt := range opts {
		if err := opt.apply(o); err != nil {
			return err
		}
	}
	return nil
}
