package mapp

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"log/slog"

	"github.com/charliego3/logger"
	"github.com/charliego3/mspp/configx"
	"github.com/charliego3/shandler"
	"github.com/gookit/goutil/strutil"

	"github.com/charliego3/mspp/grpcx"
	"github.com/charliego3/mspp/httpx"
	"github.com/charliego3/mspp/opts"
	"github.com/charliego3/mspp/service"
	"github.com/charliego3/mspp/utils"
	"github.com/soheilhy/cmux"
)

var Name string

type Application struct {
	// http server is httpx.Server using mux.Router
	http *httpx.Server

	// grpc server is grpcx.Server using grpc
	grpc *grpcx.Server

	// mux to accept http and grpc
	// if cfg.glis and cfg.hlis both nil else is nil
	mux cmux.CMux

	// Application config properties
	cfg *Config

	usingAppListener bool
}

// NewApp returns Application
func NewApp(opts ...opts.Option[Config]) *Application {
	app := &Application{}
	app.init(opts...)
	if app.cfg.onStartup != nil {
		app.cfg.onStartup(app)
	}
	return app
}

func (app *Application) fhttp(f func()) {
	if app.cfg.disableHTTP {
		return
	}

	f()
}

func (app *Application) fgrpc(f func()) {
	if app.cfg.disableGRPC {
		return
	}

	f()
}

// init handling and aggregation options
func (app *Application) init(aopts ...opts.Option[Config]) {
	app.cfg = &Config{}
	slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	for _, opt := range aopts {
		opt.Apply(app.cfg)
	}

	if app.cfg.disableGRPC && app.cfg.disableHTTP {
		slog.Error("cannot turn off HTTP and gRPC at the same time.")
		os.Exit(1)
	}

	app.usingAppListener = utils.Nils(app.cfg.glis, app.cfg.hlis) && !app.cfg.disableGRPC && !app.cfg.disableHTTP
	if app.usingAppListener {
		app.initAppListener()
		app.mux = cmux.New(app.cfg.lis)
	}

	app.fgrpc(func() {
		glogger := slog.New(shandler.NewTextHandler(shandler.WithCaller(), shandler.WithPrefix("gRPC")))
		if app.usingAppListener {
			contentType := http.CanonicalHeaderKey("content-type")
			matcher := cmux.HTTP2MatchHeaderFieldPrefixSendSettings(contentType, "application/grpc")
			app.cfg.glis = app.mux.MatchWithWriters(matcher)
		} else if app.cfg.glis == nil {
			app.cfg.glis = app.cfg.lis
			if app.cfg.glis == nil {
				app.cfg.glis = app.getDynamicListener()
			}
		}
		gopts := append(
			app.cfg.gopts,
			grpcx.WithListener(app.cfg.glis),
			grpcx.WithLogger(glogger),
		)
		app.grpc = grpcx.NewServer(gopts...)
	})

	app.fhttp(func() {
		hlogger := logger.WithPrefix("HTTP")
		if app.usingAppListener {
			app.cfg.hlis = app.mux.Match(cmux.Any())
		} else if app.cfg.hlis == nil {
			app.cfg.hlis = app.cfg.lis
			if app.cfg.hlis == nil {
				app.cfg.hlis = app.getDynamicListener()
			}
		}
		hopts := append(
			app.cfg.hopts,
			httpx.WithListener(app.cfg.hlis),
			httpx.WithLogger(hlogger),
		)
		app.http = httpx.NewServer(hopts...)
	})
}

func (app *Application) initAppListener() {
	if app.cfg.lis != nil {
		return
	}

	listener := app.getConfigListener()
	if listener == nil {
		app.cfg.lis = app.getDynamicListener()
	} else {
		app.cfg.lis = listener
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
	if app.cfg.lis == nil {
		return nil
	}
	return app.cfg.lis.Addr()
}

// RegisterService add service to http and grpc server
func (app *Application) RegisterService(services ...service.Service) {
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
				app.grpc.Logger().Error("gRPC server got an error", slog.Any("detail", err))
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
