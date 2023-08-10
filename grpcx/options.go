package grpcx

import (
	"net"

	"github.com/charliego3/logger"
	"github.com/charliego3/mspp/opts"
	"google.golang.org/grpc"
)

// WithAddr create a listener with network and address
// WithAddr and WithListener just choose one of them
func WithAddr(network, addr string) opts.Option[Server] {
	return opts.OptionFunc[Server](func(cfg *Server) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			logger.Fatal("failed to listen grpc server", "err", err)
		}
		cfg.listener = listener
	})
}

// WithListener uses the given listener
// WithListener and WithAddr just choose one of them
func WithListener(lis net.Listener) opts.Option[Server] {
	return opts.OptionFunc[Server](func(cfg *Server) {
		cfg.listener = lis
	})
}

// WithServerOption inject grpc.ServerOption to server
func WithServerOption(gsos ...grpc.ServerOption) opts.Option[Server] {
	return opts.OptionFunc[Server](func(cfg *Server) {
		cfg.srvOpts = gsos
	})
}

func WithLogger(logger logger.Logger) opts.Option[Server] {
	return opts.OptionFunc[Server](func(cfg *Server) {
		cfg.logger = logger
	})
}
