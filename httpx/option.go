package httpx

import (
	"github.com/charliego3/mspp/utility"
	"net"

	"github.com/charliego3/logger"
	"github.com/charliego3/mspp/opts"
)

// WithAddr use network and addr create a listener to serve
func WithAddr(network, addr string) utility.Option[Server] {
	return opts.OptionFunc[Server](func(cfg *Server) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			logger.Fatal("failed to listen http server", "err", err)
		}
		cfg.listener = listener
	})
}

// WithListener use this listener on Server
func WithListener(lis net.Listener) utility.Option[Server] {
	return opts.OptionFunc[Server](func(cfg *Server) {
		cfg.listener = lis
	})
}

func WithMiddleware(middles ...Middleware) utility.Option[Server] {
	return opts.OptionFunc[Server](func(cfg *Server) {
		cfg.middlewares = middles
	})
}

func WithLogger(logger logger.Logger) utility.Option[Server] {
	return opts.OptionFunc[Server](func(cfg *Server) {
		cfg.logger = logger
	})
}
