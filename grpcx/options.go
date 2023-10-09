package grpcx

import (
	"net"
	"os"

	"log/slog"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Option func(*server)

// WithAddr create a listener with network and address
// WithAddr and WithListener just choose one of them
func WithAddr(network, addr string) Option {
	return func(s *server) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			slog.Error("failed to listen gRPC server", slog.Any("err", err))
			os.Exit(1)
		}
		s.listener = listener
	}
}

// WithListener uses the given listener
// WithListener and WithAddr just choose one of them
func WithListener(lis net.Listener) Option {
	return func(s *server) {
		s.listener = lis
	}
}

// WithServerOption inject grpc.ServerOption to server
func WithServerOption(gsos ...grpc.ServerOption) Option {
	return func(s *server) {
		s.srvOpts = gsos
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *server) {
		s.logger = logger
	}
}

func WithGroup(group *errgroup.Group) Option {
	return func(s *server) {
		s.group = group
	}
}
