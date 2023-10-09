package grpcx

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/charliego3/mspp/types"
	"github.com/charliego3/mspp/utils"
	"github.com/charliego3/shandler"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type server struct {
	listener net.Listener
	server   *grpc.Server
	srvOpts  []grpc.ServerOption
	logger   *slog.Logger
	group    *errgroup.Group
}

// NewServer returns grpc server instance
func NewServer(opts ...Option) types.Server {
	s := &server{}
	for _, fn := range opts {
		fn(s)
	}
	if s.logger == nil {
		s.logger = shandler.CopyWithPrefix("gRPC")
	}
	if s.listener == nil {
		listener, err := utils.RandomTCPListener()
		if err != nil {
			s.logger.Error("failed to listen", slog.Any("err", err))
			os.Exit(1)
		}
		s.listener = listener
	}
	s.server = grpc.NewServer(s.srvOpts...)
	return s
}

func (g *server) Logger() *slog.Logger {
	return g.logger
}

// Address returns grpc listener addr
func (g *server) Address() net.Addr {
	return g.listener.Addr()
}

// RegisterService register server to grpc server
func (g *server) RegisterService(services ...types.Service) {
	for _, srv := range services {
		g.server.RegisterService(srv.Desc(), srv)
	}
}

func (g *server) Run(ctx context.Context) error {
	if g.group == nil {
		group, _ := errgroup.WithContext(ctx)
		g.group = group
	}
	g.group.Go(func() error {
		return g.server.Serve(g.listener)
	})
	g.logger.Info("listen on", slog.String("address", g.listener.Addr().String()))
	return nil
}

func (g *server) Start(ctx context.Context) error {
	err := g.Run(ctx)
	return errors.Wrap(g.group.Wait(), err.Error())
}

func (g *server) Shutdown() error {
	return nil
}
