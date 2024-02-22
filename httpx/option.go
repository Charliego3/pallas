package httpx

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/charliego3/pallas/utility"
)

type CallOption interface {
	before()
	after()
}

// WithAddr optionally specifies the TCP address for the server to listen on
func WithAddr(network, addr string) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		listener, err := net.Listen(network, addr)
		if err != nil {
			panic(err)
		}
		s.Listener = listener
	})
}

// WithListener use this listener on Server
func WithListener(lis net.Listener) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.Listener = lis
	})
}

// DisableGeneralOptions if true, passes "OPTIONS *" requests to the Handler,
// otherwise responds with 200 OK and Content-Length: 0.
func DisableGeneralOptions(disable bool) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.DisableGeneralOptionsHandler = disable
	})
}

// WithTLSConfig optionally provides a TLS configuration for use
// by ServeTLS and ListenAndServeTLS. Note that this value is
// cloned by ServeTLS and ListenAndServeTLS, so it's not
// possible to modify the configuration with methods like
// tls.Config.SetSessionTicketKeys. To use
// SetSessionTicketKeys, use Server.Serve with a TLS Listener
// instead.
func WithTLSConfig(config *tls.Config) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.TLSConfig = config
	})
}

// WithReadTimeout is the maximum duration for reading the entire
// request, including the body. A zero or negative value means
// there will be no timeout.
//
// Because ReadTimeout does not let Handlers make per-request
// decisions on each request body's acceptable deadline or
// upload rate, most users will prefer to use
// ReadHeaderTimeout. It is valid to use them both.
func WithReadTimeout(timeout time.Duration) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.ReadTimeout = timeout
	})
}

// WithReadHeaderTimeout is the amount of time allowed to read
// request headers. The connection's read deadline is reset
// after reading the headers and the Handler can decide what
// is considered too slow for the body. If ReadHeaderTimeout
// is zero, the value of ReadTimeout is used. If both are
// zero, there is no timeout.
func WithReadHeaderTimeout(timeout time.Duration) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.ReadHeaderTimeout = timeout
	})
}

// WithWriteTimeout is the maximum duration before timing out
// writes of the response. It is reset whenever a new
// request's header is read. Like ReadTimeout, it does not
// let Handlers make decisions on a per-request basis.
// A zero or negative value means there will be no timeout.
func WithWriteTimeout(timeout time.Duration) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.WriteTimeout = timeout
	})
}

// WithIdleTimeout is the maximum amount of time to wait for the
// next request when keep-alives are enabled. If IdleTimeout
// is zero, the value of ReadTimeout is used. If both are
// zero, there is no timeout.
func WithIdleTimeout(timeout time.Duration) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.IdleTimeout = timeout
	})
}

// WithMaxHeaderBytes controls the maximum number of bytes the
// server will read parsing the request header's keys and
// values, including the request line. It does not limit the
// size of the request body.
// If zero, http.DefaultMaxHeaderBytes is used.
func WithMaxHeaderBytes(max int) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.MaxHeaderBytes = max
	})
}

// WithTLSNextProto optionally specifies a function to take over
// ownership of the provided TLS connection when an ALPN
// protocol upgrade has occurred. The map key is the protocol
// name negotiated. The Handler argument should be used to
// handle HTTP requests and will initialize the Request's TLS
// and RemoteAddr if not already set. The connection is
// automatically closed when the function returns.
// If TLSNextProto is not nil, HTTP/2 support is not enabled
// automatically.
func WithTLSNextProto(proto map[string]func(*http.Server, *tls.Conn, http.Handler)) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.TLSNextProto = proto
	})
}

// WithConnState specifies an optional callback function that is
// called when a client connection changes state. See the
// ConnState type and associated constants for details.
func WithConnState(fn func(net.Conn, http.ConnState)) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.ConnState = fn
	})
}

// WithErrorLogger specifies an optional logger for errors accepting
// connections, unexpected behavior from handlers, and
// underlying FileSystem errors.
// If nil, logging is done via the log package's standard logger.
func WithErrorLogger(logger *log.Logger) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.ErrorLog = logger
	})
}

// WithBaseContext optionally specifies a function that returns
// the base context for incoming requests on this server.
// The provided Listener is the specific Listener that's
// about to start accepting requests.
// If BaseContext is nil, the default is context.Background().
// If non-nil, it must return a non-nil context.
func WithBaseContext(fn func(net.Listener) context.Context) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.BaseContext = fn
	})
}

// WithConnContext optionally specifies a function that modifies
// the context used for a new connection c. The provided ctx
// is derived from the base context and has a ServerContextKey
// value.
func WithConnContext(fn func(context.Context, net.Conn) context.Context) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.ConnContext = fn
	})
}

func WithRouterHandler(router http.Handler) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.Handler = router
	})
}

func WithMiddleware(middlewares ...Middleware) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.Router.middlewares = append(s.Router.middlewares, middlewares...)
	})
}

func WithLogger(logger *slog.Logger) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.Logger = logger
	})
}

func WithMultipartMaxSize(size int64) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.maxMultipartSize = size
	})
}

// WithStrictSlash defines the trailing slash behavior for new routes.
// The initial value is false.
//
// When true, if the route path is "/path/", accessing "/path" will perform a redirect
// to the former and vice versa. In other words, your application will always
// see the path as specified in the route.
//
// When false, if the route path is "/path", accessing "/path/" will not match
// this route and vice versa.
func WithStrictSlash(value bool) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.StrictSlash(value)
	})
}

// WithSkipClean defines the path cleaning behaviour for new routes. The initial
// value is false. Users should be careful about which routes are not cleaned
//
// When true, if the route path is "/path//to", it will remain with the double
// slash. This is helpful if you have a route like: /fetch/http://xkcd.com/534/
//
// When false, the path will be cleaned, so /fetch/http://xkcd.com/534/ will
// become /fetch/http/xkcd.com/534
func WithSkipClean(value bool) utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.SkipClean(value)
	})
}

// WithUseEncodedPath tells the router to match the encoded original path
// to the routes.
// For eg. "/path/foo%2Fbar/to" will match the path "/path/{var}/to".
//
// If not called, the router will match the unencoded path to the routes.
// For eg. "/path/foo%2Fbar/to" will match the path "/path/foo/bar/to"
func WithUseEncodedPath() utility.Option[Server] {
	return utility.OptionFunc[Server](func(s *Server) {
		s.UseEncodedPath()
	})
}
