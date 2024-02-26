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

type Middleware func(ctx *Context) (any, error)

func NewHTTPContext(req *http.Request) *Context {
	ctx := new(Context)
	ctx.Kind = KindHTTP
	ctx.Method = req.Method
	ctx.Path = req.URL.Path
	ctx.ReqHeader = Header(req.Header)
	ctx.ResHeader = make(Header)
	ctx.Payload = req.Body
	return ctx
}

func NewGRPCContext(ctx context.Context, method string, req any) *Context {
	header, _ := metadata.FromIncomingContext(ctx)
	context := new(Context)
	context.Kind = KindGRPC
	context.Path = method
	context.ReqHeader = Header(header)
	context.ResHeader = make(Header)
	context.Payload = req
	return context
}

func SetRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, req)
}

func RequestFromServerContext(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(requestKey{}).(*http.Request)
	return req, ok
}
