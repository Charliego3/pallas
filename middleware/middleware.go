package middleware

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

type Kind string

const (
	KindHTTP Kind = "HTTP"
	KindGRPC Kind = "GRPC"
)

type Context struct {
	context.Context
	Kind      Kind
	Method    string
	Path      string
	ReqHeader Header
	ResHeader Header
	Payload   any
}

type requestKey struct{}

type Handler func(ctx *Context) (any, error)

type Middleware func(next Handler) Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func (m Middleware) Append(children ...Middleware) Middleware {
	return func(next Handler) Handler {
		return Chain(children...)(m(next))
	}
}

func NewHTTPContext(req *http.Request) *Context {
	ctx := new(Context)
	ctx.Context = SetRequest(req.Context(), req)
	ctx.Kind = KindHTTP
	ctx.Method = req.Method
	ctx.Path = req.URL.Path
	ctx.ReqHeader = Header(req.Header)
	ctx.ResHeader = make(Header)
	ctx.Payload = nil
	return ctx
}

func NewGRPCContext(ctx context.Context, method string, req any) *Context {
	header, _ := metadata.FromIncomingContext(ctx)
	gctx := new(Context)
	gctx.Context = ctx
	gctx.Kind = KindGRPC
	gctx.Path = method
	gctx.ReqHeader = Header(header)
	gctx.ResHeader = make(Header)
	gctx.Payload = req
	return gctx
}

func SetRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, req)
}

func RequestFromServerContext(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(requestKey{}).(*http.Request)
	return req, ok
}
