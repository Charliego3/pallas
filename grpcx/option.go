package grpcx

import (
	"crypto/tls"
	"log/slog"
	"net"

	"github.com/charliego3/pallas/utility"

	"google.golang.org/grpc"
)

type options struct {
	serverOption  []grpc.ServerOption
	unaryInters   []grpc.UnaryServerInterceptor
	streamInters  []grpc.StreamServerInterceptor
	tlsConfig     *tls.Config
	disableHealth bool
}

func DisableHealth() utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.disableHealth = true
	})
}

func WithUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.unaryInters = interceptors
	})
}

func WithStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.streamInters = interceptors
	})
}

func WithTLS(config *tls.Config) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.tlsConfig = config
	})
}

// WithAddr create a listener with network and address
// WithAddr and WithListener just choose one of them
func WithAddr(network, addr string) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			panic(err)
		}
		s.Listener = listener
	})
}

// WithListener uses the given listener
// WithListener and WithAddr just choose one of them
func WithListener(lis net.Listener) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.Listener = lis
	})
}

// WithServerOption inject grpc.ServerOption to server
func WithServerOption(opts ...grpc.ServerOption) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.serverOption = opts
	})
}

func WithLogger(logger *slog.Logger) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.Logger = logger
	})
}
