package utility

type Option[T any] func(*T)

func Apply[T any](o *T, opts ...Option[T]) {
	for _, opt := range opts {
		opt(o)
	}
}
