package httpx

import (
	"github.com/charliego3/mspp/utility"
	"net"

	"github.com/charliego3/logger"
)

// WithAddr use network and addr create a listener to serve
func WithAddr(network, addr string) utility.Option[Server] {
	return utility.OptionFunc[Server](func(cfg *Server) error {
		listener, err := net.Listen(network, addr)
		if err != nil {
			return err
		}
		cfg.listener = listener
		return nil
	})
}

// WithListener use this listener on Server
func WithListener(lis net.Listener) utility.Option[Server] {
	return utility.InlineOpt(func(cfg *Server) {
		cfg.listener = lis
	})
}

func WithMiddleware(middles ...Middleware) utility.Option[Server] {
	return utility.InlineOpt(func(cfg *Server) {
		cfg.middlewares = middles
	})
}

func WithLogger(logger logger.Logger) utility.Option[Server] {
	return utility.InlineOpt(func(cfg *Server) {
		cfg.logger = logger
	})
}
