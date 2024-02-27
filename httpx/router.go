package httpx

import (
	"net/http"
	"path/filepath"

	"github.com/charliego3/pallas/middleware"
	"github.com/charliego3/pallas/utility"
	"github.com/gorilla/mux"
)

// This is a compilation time proposition to ensure
// that the Router is compatible with HTTP package.
var _ http.Handler = (*Router)(nil)

type Handler func(*Context) (any, error)

type ErrorEncoder func(*Context, error)

func defaultErrEncoder(c *Context, err error) {
	c.Write(map[string]any{
		"err": err.Error(),
	}, http.StatusInternalServerError)
}

type RouteWalkFunc func(method, path string)

type Router struct {
	*mux.Router
	prefix string

	// middlewares inject Middleware to route
	middlewares []middleware.Middleware

	// ene route error processor
	ene ErrorEncoder

	maxMultipartSize int64
}

func NewRouter(middlewares ...middleware.Middleware) *Router {
	r := new(Router)
	r.maxMultipartSize = 32 << 20
	r.Router = mux.NewRouter()
	r.middlewares = middlewares
	r.ene = defaultErrEncoder
	return r
}

func (r *Router) Walk(fn RouteWalkFunc) error {
	return r.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		template, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		for _, method := range methods {
			fn(method, template)
		}
		return nil
	})
}

func (r *Router) handle(method, path string, handler Handler, middlewares ...middleware.Middleware) {
	m := middleware.Chain(r.middlewares...)
	m = m.Append(middlewares...)
	next := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(w, req)
		mctx := middleware.NewHTTPContext(req)
		reply, err := m(func(c *middleware.Context) (any, error) {
			reply, err := handler(ctx)
			c.Payload = ctx.Payload
			return reply, err
		})(mctx)

		for k, v := range mctx.ResHeader {
			for _, v := range v {
				ctx.Writer.Header().Add(k, v)
			}
		}

		if err != nil {
			r.ene(ctx, err)
			return
		}

		if err = ctx.Write(reply); err != nil {
			r.ene(ctx, err)
		}
	}))
	route := r.Router.Handle(filepath.Join(r.prefix, path), next)
	if utility.NonBlank(method) {
		route.Methods(method)
	}
}

func (r *Router) Handle(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle("", path, handler, middlewares...)
}

func (r *Router) HandleFunc(path string, handler Handler, middleware ...middleware.Middleware) {
	r.handle("", path, handler, middleware...)
}

func (r *Router) HandleMethod(method, path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(method, path, handler, middlewares...)
}

func (r *Router) GET(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodGet, path, handler, middlewares...)
}

func (r *Router) POST(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodPost, path, handler, middlewares...)
}

func (r *Router) PUT(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodPut, path, handler, middlewares...)
}

func (r *Router) DELETE(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodDelete, path, handler, middlewares...)
}

func (r *Router) HEAD(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodHead, path, handler, middlewares...)
}

func (r *Router) PATCH(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodPatch, path, handler, middlewares...)
}

func (r *Router) CONNECT(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodConnect, path, handler, middlewares...)
}

func (r *Router) OPTIONS(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodOptions, path, handler, middlewares...)
}

func (r *Router) TRACE(path string, handler Handler, middlewares ...middleware.Middleware) {
	r.handle(http.MethodTrace, path, handler, middlewares...)
}

func (r *Router) Group(prefix string, middlewares ...middleware.Middleware) *Router {
	route := new(Router)
	route.prefix = filepath.Join(r.prefix, prefix)
	route.Router = r.Router
	route.middlewares = append(route.middlewares, append(r.middlewares, middlewares...)...)
	return route
}
