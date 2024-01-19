package httpx

import (
	"github.com/charliego3/mspp/utility"
	"net"
	"net/http"

	"github.com/charliego3/logger"
	"github.com/charliego3/mspp/types"
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

func NewServer(opts ...utility.Option[Server]) *Server {
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
func (h *Server) init(options []utility.Option[Server]) {
	utility.Apply(h, options...)

	if h.listener == nil {
		logger.Fatal("http server has no address specified, use WithAddr or WithListener to specify")
	}
}

func (h *Server) RegisterService(service ...types.Service) {
	h.Path("").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (h *Server) Run() error {
	h.logger.Info("serveing...")
	return http.Serve(h.listener, nil)
}
