package mapp

import (
	"context"
	"fmt"
	"github.com/charliego3/mspp/registry"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/charliego3/logger"
	"github.com/charliego3/mspp/configx"
	"github.com/charliego3/shandler"
	"github.com/gookit/goutil/strutil"

	"github.com/charliego3/mspp/grpcx"
	"github.com/charliego3/mspp/httpx"
	"github.com/charliego3/mspp/types"
	"github.com/charliego3/mspp/utility"
	"github.com/soheilhy/cmux"
)

const HealthzURL = "/debug/healthz"

// HealthzHandler is a health-check handler that returns an OK status for all
// incoming HTTP requests.
var HealthzHandler = func(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "OK")
}

var (
	// Name is the application unique name
	// this can be injected with build command `-X github.com/charliego3/mspp.Name=<app name>`
	// or using Option set value
	Name string

	// Version sepcify the application's version
	// this can be injected with build command `-X github.com/charliego3/mspp.Version=<version>`
	// or using Option set value
	Version string
)

type Application struct {
	// http server is httpx.Server using mux.Router
	http *httpx.Server

	// grpc server is grpcx.Server using grpc
	grpc types.Server

	// mux to accept http and grpc
	// if cfg.glis and cfg.hlis both nil else is nil
	mux cmux.CMux

	// Application config properties
	*options

	usingAppListener bool

	ctx      context.Context
	cancel   context.CancelCauseFunc
	servers  []types.Server
	logger   *slog.Logger
	registry registry.Registry
}

// NewApp returns Application
func NewApp(opts ...utility.Option[Application]) *Application {
	app := new(Application)
	app.logger = slog.Default()
	app.ctx, app.cancel = context.WithCancelCause(context.Background())
	utility.Apply(app, opts...)
	return app
}

func (app *Application) fhttp(f func()) {
	if app.disableHTTP {
		return
	}

	f()
}

func (app *Application) fgrpc(f func()) {
	if app.disableGRPC {
		return
	}

	f()
}

// init handling and aggregation options
func (app *Application) init(opt ...utility.Option[Application]) {
	app.options = &options{}
	slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	utility.Apply(app, opt...)

	if app.disableGRPC && app.disableHTTP {
		slog.Error("cannot turn off HTTP and gRPC at the same time.")
		os.Exit(1)
	}

	app.usingAppListener = utility.Nils(app.glis, app.hlis) && !app.disableGRPC && !app.disableHTTP
	if app.usingAppListener {
		app.initAppListener()
		app.mux = cmux.New(app.lis)
	}

	app.fgrpc(func() {
		glogger := slog.New(shandler.NewTextHandler(shandler.WithCaller(), shandler.WithPrefix("gRPC")))
		if app.usingAppListener {
			contentType := http.CanonicalHeaderKey("content-type")
			matcher := cmux.HTTP2MatchHeaderFieldPrefixSendSettings(contentType, "application/grpc")
			app.glis = app.mux.MatchWithWriters(matcher)
		} else if app.glis == nil {
			app.glis = app.lis
			if app.glis == nil {
				app.glis = app.getDynamicListener()
			}
		}
		gopts := []utility.Option[grpcx.Server]{
			grpcx.WithListener(app.glis),
			grpcx.WithLogger(glogger),
		}
		app.grpc = grpcx.NewServer(gopts...)
	})

	app.fhttp(func() {
		hlogger := logger.WithPrefix("HTTP")
		if app.usingAppListener {
			app.hlis = app.mux.Match(cmux.Any())
		} else if app.hlis == nil {
			app.hlis = app.lis
			if app.hlis == nil {
				app.hlis = app.getDynamicListener()
			}
		}
		hopts := append(
			app.hopts,
			httpx.WithListener(app.hlis),
			httpx.WithLogger(hlogger),
		)
		app.http = httpx.NewServer(hopts...)
	})
}

func (app *Application) initAppListener() {
	if app.lis != nil {
		return
	}

	listener := app.getConfigListener()
	if listener == nil {
		app.lis = app.getDynamicListener()
	} else {
		app.lis = listener
	}
}

func (app *Application) getConfigListener() net.Listener {
	cfg, err := configx.Fetch[configx.App]()
	if err == nil && strutil.IsNotBlank(cfg.Address) {
		if strutil.IsBlank(cfg.Network) {
			cfg.Network = "tcp"
		}
		_, _, err = net.SplitHostPort(cfg.Address)
		if err != nil {
			slog.Error("failed listen application",
				slog.String("network", cfg.Network),
				slog.String("address", cfg.Address),
				slog.Any("reason", err),
			)
			os.Exit(1)
		}

		listener, err := net.Listen(cfg.Network, cfg.Address)
		if err != nil {
			slog.Error("failed listen application",
				slog.String("network", cfg.Network),
				slog.String("address", cfg.Address),
				slog.Any("reason", err),
			)
			os.Exit(1)
		}
		return listener
	}
	return nil
}

// getDynamicListener if app without any listener specifies then create a dynamic listener
func (app *Application) getDynamicListener() net.Listener {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		slog.Error("failed to listen", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Warn("No address is specified so dynamic addresses are used", slog.String("address", listener.Addr().String()))
	return listener
}

// Address returns application listen address
// this address is http and grpc both
func (app *Application) Address() net.Addr {
	if app.lis == nil {
		return nil
	}
	return app.lis.Addr()
}

// RegisterService add service to http and grpc server
func (app *Application) RegisterService(services ...types.Service) {
	app.fhttp(func() {
		app.http.RegisterService(services...)
	})
	app.fgrpc(func() {
		app.grpc.RegisterService(services...)
	})
}

// Run start the server until terminate
func (app *Application) Run(ctx context.Context) (err error) {
	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)
	group.Go(func() error {
		app.fgrpc(func() {
			err := app.grpc.Run(ctx)
			if err != nil {
				//app.grpc.Logger.Error("gRPC server got an error", slog.Any("detail", err))
			}
		})
		return err
	})

	go app.fhttp(func() {
		err := app.http.Run()
		if err != nil {
			app.http.Logger().Fatal("HTTP server got an error", err)
		}
	})

	if app.usingAppListener {
		go func() {
			err := app.mux.Serve()
			if err != nil {
				slog.Error("Application got an error", slog.Any("err", err))
				os.Exit(1)
			}
		}()
	}

	addr := app.Address()
	if addr != nil {
		slog.Info("listening", slog.String("address", addr.String()))
	}

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGINT)
	<-stopper

	slog.Info("terminated.")
	return err
}

func (app *Application) Shutdown() {

}
