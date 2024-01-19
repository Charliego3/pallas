package grpcx

import (
	"github.com/charliego3/mspp/utility"
	"log/slog"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type options struct {
	serverOpts   []grpc.ServerOption
	unaryInters  []grpc.UnaryServerInterceptor
	streamInters []grpc.StreamServerInterceptor
}

// WithAddr create a listener with network and address
// WithAddr and WithListener just choose one of them
func WithAddr(network, addr string) utility.Option[Server] {
	return func(s *Server) {
		s.base.Network = network
		s.base.Addr = addr
	}
}

// WithListener uses the given listener
// WithListener and WithAddr just choose one of them
func WithListener(lis net.Listener) utility.Option[Server] {
	return func(s *Server) {
		s.base.Listener = lis
	}
}

// WithServerOption inject grpc.ServerOption to server
func WithServerOption(opts ...grpc.ServerOption) utility.Option[Server] {
	return func(s *Server) {
		s.serverOpts = opts
	}
}

func WithLogger(logger *slog.Logger) utility.Option[Server] {
	return func(s *Server) {
		s.base.Logger = logger
	}
}

func WithGroup(group *errgroup.Group) utility.Option[Server] {
	return func(s *Server) {
		s.group = group
	}
}
