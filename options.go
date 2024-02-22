package pallas

import (
	"context"
	"log/slog"
	"net"

	"github.com/charliego3/pallas/grpcx"
	"github.com/charliego3/pallas/utility"

	"github.com/charliego3/pallas/httpx"
)

type StartInterceptor func(context.Context) (context.Context, error)

type options struct {
	// http and grpc server both listen this address
	// when gopts or hopts has not sepcity listener
	listener net.Listener

	// gopts is grpcx.Server options
	gopts []utility.Option[grpcx.Server]

	// middles accept http server Middleware
	hopts []utility.Option[httpx.Server]

	middlewares []Middleware

	// onStartup run on Applition after init
	onStartup   func(*Application) error
	beforeStart StartInterceptor
	afterStart  StartInterceptor

	beforeShutdown func(context.Context) error
	afterShutdown  func(context.Context) error
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
