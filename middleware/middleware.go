package middleware

import (
	"context"
	"net/http"
)

type Kind uint

const (
	KindHTTP Kind = iota
	KindGRPC
)

type Context struct {
	Kind       Kind
	Method     string
	Path       string
	ReqHeader  Header
	ReplyHader Header
	Payload    any
}

type requestKey struct{}

type Middleware func(ctx *Context) (any, error)

func SetRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, req)
}

func RequestFromServerContext(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(requestKey{}).(*http.Request)
	return req, ok
}
