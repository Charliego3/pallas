package grpcx

import (
	"context"
	"github.com/charliego3/mspp/types"
	"github.com/charliego3/mspp/utility"
	"github.com/charliego3/shandler"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

var NoListener = errors.New("gRPC server")

type Server struct {
	*options
	base   *types.BaseServer
	server *grpc.Server
	health grpc_health_v1.HealthServer
	group  *errgroup.Group
}

// NewServer returns grpc server instance
func NewServer(opts ...utility.Option[Server]) *Server {
	s := new(Server)
	s.base = new(types.BaseServer)
	s.options = new(options)
	s.health = health.NewServer()
	utility.Apply(s, opts...)
	if s.base.Logger == nil {
		s.base.Logger = shandler.CopyWithPrefix("gRPC")
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(append(
			[]grpc.UnaryServerInterceptor{unaryInterceptor},
			s.unaryInters...,
		)...),
		grpc.ChainStreamInterceptor(append(
			[]grpc.StreamServerInterceptor{streamInterceptor},
			s.streamInters...,
		)...),
	}
	if s.tlsConfig != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(s.tlsConfig)))
	}
	if len(s.serverOption) > 0 {
		grpcOpts = append(grpcOpts, s.serverOption...)
	}
	s.server = grpc.NewServer(s.serverOption...)
	if !s.disableHealth {
		grpc_health_v1.RegisterHealthServer(s.server, s.health)
	}
	reflection.Register(s.server)
	return s
}

// Listener returns grpc listener
func (g *Server) Listener() net.Listener {
	return g.base.Listener
}

// RegisterService register server to grpc server
func (g *Server) RegisterService(services ...types.Service) {
	for _, srv := range services {
		g.server.RegisterService(srv.Desc(), srv)
	}
}

func (g *Server) Run(ctx context.Context) error {
	if !g.base.HasListener() {
		return NoListener
	}

	if g.base.Listener == nil {
		listener, err := net.Listen(g.base.Network, g.base.Addr)
		if err != nil {
			return err
		}

		g.base.Listener = listener
	}
	return nil
}

func (g *Server) Start(ctx context.Context) error {
	err := g.Run(ctx)
	return errors.Wrap(g.group.Wait(), err.Error())
}

func (g *Server) Shutdown(ctx context.Context) error {
	return nil
}
