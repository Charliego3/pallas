package grpcx

import (
	"crypto/tls"
	"github.com/charliego3/mspp/utility"
	"log/slog"
	"net"

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
	return utility.InlineOpt(func(s *Server) {
		s.disableHealth = true
	})
}

func WithUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) utility.Option[Server] {
	return utility.InlineOpt(func(s *Server) {
		s.unaryInters = interceptors
	})
}

func WithStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) utility.Option[Server] {
	return utility.InlineOpt(func(s *Server) {
		s.streamInters = interceptors
	})
}

func WithTLS(config *tls.Config) utility.Option[Server] {
	return utility.InlineOpt(func(s *Server) {
		s.tlsConfig = config
	})
}

// WithAddr create a listener with network and address
// WithAddr and WithListener just choose one of them
func WithAddr(network, addr string) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) error {
		if listener, err := net.Listen(network, addr); err != nil {
			return err
		} else {
			s.Listener = listener
		}
		return nil
	})
}

// WithListener uses the given listener
// WithListener and WithAddr just choose one of them
func WithListener(lis net.Listener) utility.Option[Server] {
	return utility.InlineOpt(func(s *Server) {
		s.Listener = lis
	})
}

// WithServerOption inject grpc.ServerOption to server
func WithServerOption(opts ...grpc.ServerOption) utility.Option[Server] {
	return utility.InlineOpt(func(s *Server) {
		s.serverOption = opts
	})
}

func WithLogger(logger *slog.Logger) utility.Option[Server] {
	return utility.InlineOpt(func(s *Server) {
		s.Logger = logger
	})
}
