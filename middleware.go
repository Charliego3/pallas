package pallas

import (
	// "context"
	// "net"
	"sync"

	"github.com/charliego3/pallas/httpx"
	"google.golang.org/grpc"
)

// type GrpcUnaryMiddleware interface {
// 	Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
// }

// type GrpcStreamMiddleware interface {
// 	Stream(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
// }

// type Middleware interface {
// 	HTTP(next httpx.Handler) httpx.Handler
// 	GrpcUnaryMiddleware
// 	GrpcStreamMiddleware
// }

// var _ Middleware = (*loggingMiddleware)(nil)

// type loggingMiddleware struct{}

// func (*loggingMiddleware) HTTP(next httpx.Handler) httpx.Handler {
// 	return next
// }

// func (*loggingMiddleware) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	return handler(ctx, req)
// }

// func (*loggingMiddleware) Stream(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 	return handler(srv, ss)
// }
//

type MiddlewareKind uint

const (
	MiddlewareKindGrcp MiddlewareKind = iota
	MiddlewareKindHttp
)

type MiddleContext struct {
	*httpx.Context
	kind    MiddlewareKind
	req     any // Grpc Unary request object
	info    any // Grpc Unary/Stream server info
	handler any // Grpc Unary/Stream server handler

	srv any // Grpc Stream server
	ss  grpc.ServerStream
}

func (c MiddleContext) RequestObj() any {
	return c.req
}

func (c MiddleContext) UnaryServerInfo() *grpc.UnaryServerInfo {
	if si, ok := c.info.(*grpc.UnaryServerInfo); ok {
		return si
	}
	return nil
}

func (c MiddleContext) StreamServerInfo() *grpc.StreamServerInfo {
	if si, ok := c.info.(*grpc.StreamServerInfo); ok {
		return si
	}
	return nil
}

func (c MiddleContext) UnaryHandler() grpc.UnaryHandler {
	if uh, ok := c.handler.(grpc.UnaryHandler); ok {
		return uh
	}
	return nil
}

func (c MiddleContext) StreamHandler() grpc.StreamHandler {
	if sh, ok := c.handler.(grpc.StreamHandler); ok {
		return sh
	}
	return nil
}

func (c MiddleContext) StreamServer() any {
	return c.srv
}

func (c MiddleContext) ServerStream() grpc.ServerStream {
	return c.ss
}

func (c MiddleContext) Kind() MiddlewareKind {
	return c.kind
}

type Middleware func(*MiddleContext) error

var ctxpool = &sync.Pool{
	New: func() any {
		return new(MiddleContext)
	},
}

func LoggingMiddleware(c *MiddleContext) error {
	if c.Kind() == MiddlewareKindHttp {

	}
	return nil
}

var _ Middleware = LoggingMiddleware
