package mapp

import (
	"context"
	"github.com/charliego3/mspp/types"
	"github.com/charliego3/mspp/utility"
	"net"
	"os"

	"github.com/charliego3/mspp/httpx"
	"google.golang.org/grpc"
	"log/slog"
)

type StartInterceptor func(context.Context) (context.Context, error)

type options struct {
	// http and grpc server both listen this address
	lis net.Listener

	// http server to serve this address
	// if lis and hlis both nil, then hlis using dynamic address
	hlis net.Listener

	// grpc server to serve this address
	// if lis and glis both nil, then glis using dynamic address
	glis net.Listener

	// disableHTTP only serve grpc server
	disableHTTP bool

	// disableGRPC only serve http server
	disableGRPC bool

	// onStartup run on Applition after init
	onStartup func(*Application) error

	// gopts is grpcx.Server options
	//gopts []grpcx.Option

	// middles accept http server Middleware
	hopts []utility.Option[httpx.Server]

	beforeStart StartInterceptor
	afterStart  StartInterceptor

	beforeShutdown func(context.Context) error
	afterShutdown  func(context.Context) error
}

func WithName(name string) utility.Option[Application] {
	return func(*Application) {
		Name = name
	}
}

func WithVersion(version string) utility.Option[Application] {
	return func(*Application) {
		Version = version
	}
}

func WithLogger(logger *slog.Logger) utility.Option[Application] {
	return func(app *Application) {
		app.logger = logger
	}
}

func WithServer(sv ...types.Server) utility.Option[Application] {
	return func(app *Application) {
		app.servers = append(app.servers, sv...)
	}
}

func WithBeforeStart(handler StartInterceptor) utility.Option[Application] {
	return func(app *Application) {
		app.beforeStart = handler
	}
}

func DisableHTTP() utility.Option[options] {
	return func(cfg *options) {
		cfg.disableHTTP = true
	}
}

func DisableGRPC() utility.Option[options] {
	return func(cfg *options) {
		cfg.disableGRPC = true
	}
}

func WithHttpOpts(hopts ...utility.Option[httpx.Server]) utility.Option[options] {
	return func(cfg *options) {
		cfg.hopts = hopts
	}
}

func OnStartup(fn func(*Application) error) utility.Option[options] {
	return func(cfg *options) {
		cfg.onStartup = fn
	}
}

// WithGrpcServerOpts accept grpc server options
func WithGrpcServerOpts(gopts ...grpc.ServerOption) utility.Option[options] {
	return func(cfg *options) {
		//cfg.gopts = append(cfg.gopts, grpcx.WithServerOption(gopts...))
	}
}

// WithAddr served http and grpc on same address
func WithAddr(network, addr string) utility.Option[options] {
	return func(cfg *options) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			slog.Error("failed to listen app", slog.Any("err", err))
			os.Exit(1)
		}
		cfg.lis = listener
	}
}

// WithTCPAddr is WithAddr alias but network using tcp
func WithTCPAddr(addr string) utility.Option[options] {
	return WithAddr("tcp", addr)
}

// WithListener served http and grpc on same address
func WithListener(lis net.Listener) utility.Option[options] {
	return func(cfg *options) {
		cfg.lis = lis
	}
}

// WithHttpAddr expected http server listen address using tcp network
func WithHttpAddr(network, addr string) utility.Option[options] {
	return func(cfg *options) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			slog.Error("failed to listen http server with app", slog.Any("err", err))
			os.Exit(1)
		}
		cfg.hlis = listener
	}
}

// WithHttpListener served http server listener
func WithHttpListener(lis net.Listener) utility.Option[options] {
	return func(cfg *options) {
		cfg.hlis = lis
	}
}

// WithGrpcAddr served grpc server on address
func WithGrpcAddr(network, addr string) utility.Option[options] {
	return func(cfg *options) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			slog.Error("failed to listen grpc server with app", slog.Any("err", err))
			os.Exit(1)
		}
		cfg.glis = listener
	}
}

// WithGrpcListener served grpc server listener
func WithGrpcListener(lis net.Listener) utility.Option[options] {
	return func(cfg *options) {
		cfg.glis = lis
	}
}
