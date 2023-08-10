package configx

import (
	"fmt"
	"sync"
)

type Fetcher[T any] interface {
	Fetch() (T, error)
}

type CachedFetcher[T any] struct {
	fetcher Fetcher[T]
	payload T
	once    sync.Once
}

func RegisterCachedFetcher[T any](fetcher Fetcher[T]) {
	RegisterFetcher[T](&CachedFetcher[T]{fetcher: fetcher})
}

func (f *CachedFetcher[T]) Fetch() (T, error) {
	var err error
	f.once.Do(func() {
		f.payload, err = f.fetcher.Fetch()
	})
	return f.payload, err
}

var fetchers = make([]any, 0)

func RegisterFetcher[T any](fetcher Fetcher[T]) {
	fetchers = append(fetchers, fetcher)
}

func Fetch[T any]() (T, error) {
	return getConfig[T]()
}

func Must[T any]() T {
	ins, _ := Fetch[T]()
	return ins
}

func getFetcher[T any]() (Fetcher[T], bool) {
	var t T
	if (any)(t) == nil {
		return nil, false
	}

	for i := len(fetchers) - 1; i >= 0; i-- {
		if f, ok := fetchers[i].(Fetcher[T]); ok {
			return f, true
		}
	}
	return nil, false
}

func getConfig[T any]() (t T, ok error) {
	if f, ok := getFetcher[T](); ok {
		return f.Fetch()
	}
	return t, fmt.Errorf("config not found: %T", t)
}
