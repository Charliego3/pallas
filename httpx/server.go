package httpx

import (
	"net"
	"net/http"

	"github.com/charliego3/logger"
	"github.com/charliego3/mspp/opts"
	"github.com/charliego3/mspp/service"
	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router

	// this listener will be served
	listener net.Listener

	// http server middleware
	middlewares []Middleware

	logger logger.Logger
}

type Middleware = mux.MiddlewareFunc

func NewServer(opts ...opts.Option[Server]) *Server {
	h := &Server{
		Router: mux.NewRouter(),
	}

	h.init(opts)
	h.Use(h.middlewares...)
	return h
}

func (h *Server) Logger() logger.Logger {
	return h.logger
}

// init accept options to Server
func (h *Server) init(opts []opts.Option[Server]) {
	for _, opt := range opts {
		opt.Apply(h)
	}

	if h.listener == nil {
		logger.Fatal("http server has no address specified, use WithAddr or WithListener to specify")
	}
}

func (h *Server) RegisterService(service ...service.Service) {
	h.Path("").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (h *Server) Run() error {
	h.logger.Info("serveing...")
	return http.Serve(h.listener, nil)
}
