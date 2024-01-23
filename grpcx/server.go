package grpcx

import (
	"context"
	"github.com/charliego3/mspp/types"
	"github.com/charliego3/mspp/utility"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log/slog"
)

var (
	_ types.Server = (*Server)(nil)

	NoListener = errors.New("[gRPC] server not bind listener")
)

type Server struct {
	*options
	*types.BaseServer
	server *grpc.Server
	health *health.Server
	ctx    context.Context
	err    error
}

// NewServer returns grpc server instance
func NewServer(opts ...utility.Option[Server]) *Server {
	s := new(Server)
	s.BaseServer = new(types.BaseServer)
	s.options = new(options)
	s.health = health.NewServer()
	s.err = utility.Apply(s, opts...)
	if s.Logger == nil {
		s.Logger = slog.Default()
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
	s.server = grpc.NewServer(grpcOpts...)
	if !s.disableHealth {
		grpc_health_v1.RegisterHealthServer(s.server, s.health)
	}
	reflection.Register(s.server)
	return s
}

// RegisterService register server to grpc server
func (g *Server) RegisterService(services ...types.Service) {
	if g.err != nil {
		return
	}

	for _, srv := range services {
		desc := srv.Desc()
		g.server.RegisterService(&desc, srv)
	}
}

func (g *Server) Run(ctx context.Context) error {
	if g.err != nil {
		return g.err
	}

	if g.Listener == nil {
		return NoListener
	}

	g.ctx = ctx
	g.health.Resume()
	g.Logger.Info("[gRPC] server listening on", slog.String("address", g.Listener.Addr().String()))
	return g.server.Serve(g.Listener)
}

func (g *Server) Shutdown(context.Context) error {
	g.health.Shutdown()
	g.server.GracefulStop()
	return nil
}
