package middleware

import (
	"strings"
)

type Header map[string][]string

func (h Header) Get(key string) string {
	vals := h[strings.ToLower(key)]
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (h Header) Set(key string, vals ...string) {
	h[strings.ToLower(key)] = vals
}

func (h Header) Add(key string, vals ...string) {
	if len(vals) == 0 {
		return
	}

	key = strings.ToLower(key)
	h[key] = append(h[key], vals...)
}

func (h Header) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

func (h Header) Values(key string) []string {
	return h[strings.ToLower(key)]
}

func (h Header) Copy() Header {
	header := make(Header, len(h))
	for k, vals := range h {
		header[k] = copyOf(vals)
	}
	return header
}

func (h Header) Len() int {
	return len(h)
}

func copyOf(v []string) []string {
	vals := make([]string, len(v))
	copy(vals, v)
	return vals
}
