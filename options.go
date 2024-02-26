package pallas

import (
	"context"
	"log/slog"
	"net"

	"github.com/charliego3/pallas/grpcx"
	"github.com/charliego3/pallas/middleware"
	"github.com/charliego3/pallas/utility"
	"github.com/soheilhy/cmux"

	"github.com/charliego3/pallas/httpx"
)

type StartInterceptor func(context.Context) (context.Context, error)

type options struct {
	// http and grpc server both listen this address
	// when gopts or hopts has not sepcity listener
	listener net.Listener

	// grpcMatcher match the grpc request on same listener
	// default using header content-type: application/grpc
	grpcMatcher cmux.MatchWriter

	// gopts is grpcx.Server options
	gopts []utility.Option[grpcx.Server]

	// middles accept http server Middleware
	hopts []utility.Option[httpx.Server]

	// onStartup run on Applition after init
	onStartup   func(*Application) error
	beforeStart StartInterceptor
	afterStart  StartInterceptor

	beforeShutdown func(context.Context) error
	afterShutdown  func(context.Context) error
}

// WithGrpcMatcher custom grpc dispatcher
func WithGrpcMatcher(matcher cmux.MatchWriter) utility.Option[Application] {
	return utility.OptionFunc[Application](func(app *Application) {
		app.grpcMatcher = matcher
	})
}

func WithMiddleware(middlewares ...middleware.Middleware) utility.Option[Application] {
	return utility.OptionFunc[Application](func(app *Application) {
		app.hopts = append(app.hopts, httpx.WithMiddleware(middlewares...))
		app.gopts = append(app.gopts, grpcx.WithMiddleware(middlewares...))
	})
}

func WithLogger(logger *slog.Logger) utility.Option[Application] {
	return utility.OptionFunc[Application](func(app *Application) {
		app.logger = logger
	})
}

func WithBeforeStart(handler StartInterceptor) utility.Option[Application] {
	return utility.OptionFunc[Application](func(app *Application) {
		app.beforeStart = handler
	})
}

func WithHttpOpts(hopts ...utility.Option[httpx.Server]) utility.Option[Application] {
	return utility.OptionFunc[Application](func(cfg *Application) {
		cfg.hopts = append(cfg.hopts, hopts...)
	})
}

func OnStartup(fn func(app *Application) error) utility.Option[Application] {
	return utility.OptionFunc[Application](func(cfg *Application) {
		cfg.onStartup = fn
	})
}

// WithGrpcOpts accept grpc server options
func WithGrpcOpts(gopts ...utility.Option[grpcx.Server]) utility.Option[Application] {
	return utility.OptionFunc[Application](func(cfg *Application) {
		cfg.gopts = append(cfg.gopts, gopts...)
	})
}

// WithAddr served http and grpc on same address
func WithAddr(network, addr string) utility.Option[Application] {
	return utility.OptionFunc[Application](func(cfg *Application) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			panic(err)
		}
		cfg.listener = listener
	})
}

// WithTCPAddr is WithAddr alias but network using tcp
func WithTCPAddr(addr string) utility.Option[Application] {
	return WithAddr("tcp", addr)
}

// WithListener served http and grpc on same address
func WithListener(listener net.Listener) utility.Option[Application] {
	return utility.OptionFunc[Application](func(cfg *Application) {
		cfg.listener = listener
	})
}

// WithDefaultCodecType is http server defalt Codec type name
// encoding/json.Type or encoding/xml.Type or custom register Codec type
func WithDefaultCodecType(typename string) utility.Option[Application] {
	return utility.OptionFunc[Application](func(cfg *Application) {
		httpx.SetDefaultCodeType(typename)
	})
}
