package pallas

import (
	"context"
	"errors"
	"fmt"
	"github.com/charliego3/pallas/registry"
	"golang.org/x/sync/errgroup"
	"net"

	"log/slog"

	"github.com/charliego3/pallas/configx"

	"github.com/charliego3/pallas/grpcx"
	"github.com/charliego3/pallas/httpx"
	"github.com/charliego3/pallas/types"
	"github.com/charliego3/pallas/utility"
	"github.com/soheilhy/cmux"
)

var (
	// Name is the application unique name
	// this can be injected with build command `-X github.com/charliego3/pallas.Name=<app name>`
	// or using Option set value
	Name string

	// Version sepcify the application's version
	// this can be injected with build command `-X github.com/charliego3/pallas.Version=<version>`
	// or using Option set value
	Version string
)

type Application struct {
	// http server is httpx.Server handler http request
	http *httpx.Server

	// grpc server is grpcx.Server handler grpc request
	grpc *grpcx.Server

	// mux to accept http and grpc
	// if http.Listener and grpc.Listener both nil will be using else is nil
	mux cmux.CMux

	// Application config properties
	*options

	ctx      context.Context
	cancel   context.CancelCauseFunc
	logger   *slog.Logger
	registry registry.Registry
}

// NewApp returns Application
func NewApp(opts ...utility.Option[Application]) *Application {
	app := new(Application)
	app.logger = slog.Default()
	app.options = new(options)
	utility.Apply(app, opts...)

	app.http = httpx.NewServer(app.hopts...)
	app.grpc = grpcx.NewServer(app.gopts...)
	if utility.Nils(app.http.Listener, app.grpc.Listener) {
		app.mux = cmux.New(app.getListener())
		matcher := cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc")
		app.grpc.Listener = app.mux.MatchWithWriters(matcher)
		app.http.Listener = app.mux.Match(cmux.Any())
	} else {
		listener, err := utility.RandomTCPListener()
		if err != nil {
			panic(fmt.Sprintf("bind dynamic listener paniced: %v", err))
		}
		app.grpc.Listener = utility.DObj(app.grpc.Listener, listener)
		app.http.Listener = utility.DObj(app.http.Listener, listener)
	}
	return app
}

// getListener first using app options listener
// otherwise using listener from config
// last using a dynamic listener
func (app *Application) getListener() (listener net.Listener) {
	if app.listener != nil {
		return app.listener
	}

	cfg, err := configx.Fetch[configx.App]()
	if err != nil && !errors.Is(err, configx.ErrNotFound) {
		panic(fmt.Sprintf("fetch App config paniced: %v", err))
	}

	listener, err = net.Listen(
		utility.DString(cfg.Network, "tcp"),
		utility.DString(cfg.Address, ":0"),
	)
	if err != nil {
		panic(fmt.Sprintf("create App listener paniced: %v", err))
	}
	return listener
}

// Address returns application listen address
// this address is http and grpc both
func (app *Application) Address() net.Addr {
	if app.listener == nil {
		return nil
	}
	return app.listener.Addr()
}

// RegisterService add service to http and grpc server
func (app *Application) RegisterService(services ...types.Service) {
	app.http.RegisterService(services...)
	app.grpc.RegisterService(services...)
}

// Run start the server until terminate
func (app *Application) Run(ctx context.Context) (err error) {
	group, ctx := errgroup.WithContext(ctx)
	app.ctx = ctx
	group.Go(func() error {
		return app.http.Run(ctx)
	})
	group.Go(func() error {
		return app.grpc.Run(ctx)
	})
	if app.mux != nil {
		group.Go(app.mux.Serve)
	}
	app.logger.Info("started.")
	return group.Wait()
}

func (app *Application) Shutdown() {

}
