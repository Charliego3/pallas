package grpcx

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/charliego3/mspp/types"
	"github.com/charliego3/mspp/utility"
	"github.com/charliego3/shandler"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	*options
	base   *types.BaseServer
	server *grpc.Server
	group  *errgroup.Group
}

// NewServer returns grpc server instance
func NewServer(opts ...utility.Option[Server]) *Server {
	s := new(Server)
	s.base = new(types.BaseServer)
	s.options = new(options)
	utility.Apply(s, opts...)
	if s.base.Logger == nil {
		s.base.Logger = shandler.CopyWithPrefix("gRPC")
	}
	if s.base.Listener == nil {
		listener, err := utility.RandomTCPListener()
		if err != nil {
			s.base.Logger.Error("failed to listen", slog.Any("err", err))
			os.Exit(1)
		}
		s.base.Listener = listener
	}
	s.server = grpc.NewServer(s.serverOpts...)
	return s
}

// Address returns grpc listener addr
func (g *Server) Address() net.Addr {
	return g.base.Listener.Addr()
}

// RegisterService register server to grpc server
func (g *Server) RegisterService(services ...types.Service) {
	for _, srv := range services {
		g.server.RegisterService(srv.Desc(), srv)
	}
}

func (g *Server) Run(ctx context.Context) error {
	if g.group == nil {
		group, _ := errgroup.WithContext(ctx)
		g.group = group
	}
	g.group.Go(func() error {
		return g.server.Serve(g.base.Listener)
	})
	g.base.Logger.Info("listen on", slog.String("address", g.base.Listener.Addr().String()))
	return nil
}

func (g *Server) Start(ctx context.Context) error {
	err := g.Run(ctx)
	return errors.Wrap(g.group.Wait(), err.Error())
}

func (g *Server) Shutdown(ctx context.Context) error {
	return nil
}
