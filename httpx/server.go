package httpx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/charliego3/pallas/types"
	"github.com/charliego3/pallas/utility"
)

const HealthzURI = "/debug/healthz"

// HealthzHandler is a health-check handler that returns an OK status for all
// incoming HTTP requests.
var HealthzHandler = func(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "OK")
}

type Server struct {
	*options
	*types.BaseServer
	*http.Server
	*Router
}

func NewServer(opts ...utility.Option[Server]) *Server {
	h := new(Server)
	h.Router = NewRouter()
	h.Server = new(http.Server)
	h.options = newDefauleOpts()
	h.BaseServer = types.NewDefaultBaseServer()
	utility.Apply(h, opts...)
	return h
}

func (h *Server) RegisterService(service ...types.Service) {
	if h.Handler != nil {
		return
	}

	for _, serv := range service {
		h.registerService(serv)
	}
}

func (h *Server) registerService(srv types.Service) {
	for _, m := range srv.Desc().Http.Methods {
		handler := m.Handler(srv)
		if handler == nil {
			panic("nil handler cannot register on httpx.Server")
		}
		if hd, ok := handler.(Handler); !ok {
			panic(fmt.Sprintf("%T handler cannot register, expect: httpx.Handler", handler))
		} else {
			h.handle(m.Method, m.Template, hd)
		}
	}
}

func (h *Server) Run(ctx context.Context) error {
	if h.Listener == nil {
		return errors.New("[HTTP] not bind listener")
	}

	if h.Handler == nil {
		h.Handler = h.Router
	}
	h.Logger.Info("[HTTP] listening on", slog.String("address", h.Listener.Addr().String()))
	if h.TLSConfig != nil {
		return h.ServeTLS(h.Listener, "", "")
	}
	return h.Serve(h.Listener)
}

func (h *Server) Shutdown(ctx context.Context) error {
	return h.Server.Shutdown(ctx)
}
