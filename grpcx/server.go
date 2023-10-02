package grpcx

import (
	"context"
	"net"
	"os"

	"github.com/charliego3/mspp/service"
	"github.com/charliego3/mspp/utils"
	"github.com/charliego3/shandler"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	listener net.Listener
	server   *grpc.Server
	srvOpts  []grpc.ServerOption
	logger   *slog.Logger
	group    *errgroup.Group
}

// NewServer returns grpc server instance
func NewServer(opts ...Option) *Server {
	s := &Server{}
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

func (g *Server) Logger() *slog.Logger {
	return g.logger
}

// Address returns grpc listener addr
func (g *Server) Address() net.Addr {
	return g.listener.Addr()
}

// RegisterService register server to grpc servser
func (g *Server) RegisterService(services ...service.Service) {
	for _, srv := range services {
		g.server.RegisterService(srv.ServiceDesc(), srv)
	}
}

func (g *Server) Run(ctx context.Context) error {
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

func (g *Server) Wait() error {
	return g.group.Wait()
}
