package types

type Optional[T any] struct{ v T }

func OptionalOf[T any](v T) Optional[T] {
	return Optional[T]{v}
}

func OptionalEmpty[T any]() Optional[T] {
	return OptionalEmpty[T]()
}

func OptionalMap[T, U any](opt Optional[T], mapper func(T) U) Optional[U] {
	if opt.IsEmpty() {
		return OptionalEmpty[U]()
	}
	return OptionalOf[U](mapper(opt.v))
}

func (t Optional[T]) IfPresent(f func(T)) {
	if t.IsEmpty() {
		return
	}

	f(t.v)
}

func (t Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if t.IsEmpty() {
		return t
	}
	if predicate(t.v) {
		return t
	}
	return OptionalEmpty[T]()
}

func (t Optional[T]) IsEmpty() bool {
	return (any)(t.v) == nil
}

func (t Optional[T]) IsPresent() bool {
	return !t.IsEmpty()
}

func (t Optional[T]) Else(other T) T {
	if t.IsPresent() {
		return t.v
	}
	return other
}
