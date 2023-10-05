package grpcx

import (
	"net"

	"github.com/charliego3/logger"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log/slog"
)

type Option func(*Server)

// WithAddr create a listener with network and address
// WithAddr and WithListener just choose one of them
func WithAddr(network, addr string) Option {
	return func(s *Server) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			logger.Fatal("failed to listen grpc server", "err", err)
		}
		s.listener = listener
	}
}

// WithListener uses the given listener
// WithListener and WithAddr just choose one of them
func WithListener(lis net.Listener) Option {
	return func(s *Server) {
		s.listener = lis
	}
}

// WithServerOption inject grpc.ServerOption to server
func WithServerOption(gsos ...grpc.ServerOption) Option {
	return func(s *Server) {
		s.srvOpts = gsos
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithGroup(group *errgroup.Group) Option {
	return func(s *Server) {
		s.group = group
	}
}
