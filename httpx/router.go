package httpx

import (
	"net/http"
	"path/filepath"

	"github.com/charliego3/pallas/utility"
	"github.com/gorilla/mux"
)

// This is a compilation time proposition to ensure
// that the Router is compatible with HTTP package.
var _ http.Handler = (*Router)(nil)

type Handler interface {
	Serve(c *Context) error
}

type HandlerFunc func(c *Context) error

func (f HandlerFunc) Serve(c *Context) error {
	return f(c)
}

type ErrorEncoder func(*Context, error)

func defaultErrEncoder(c *Context, err error) {
	c.Write(err)
}

type RouteWalkFunc func(method, path string)

type Router struct {
	*mux.Router
	prefix string

	// middlewares inject Middleware to route
	middlewares []Middleware

	// ene route error processor
	ene ErrorEncoder

	maxMultipartSize int64
}

func NewRouter(middlewares ...Middleware) *Router {
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

func (r *Router) mergeMiddlewares(handler Handler, middlewares ...Middleware) Handler {
	middlewares = append(r.middlewares, middlewares...)
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (r *Router) handle(method, path string, handler Handler, middlewares ...Middleware) {
	next := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler = r.mergeMiddlewares(handler, middlewares...)
		ctx := NewContext(w, req)
		if err := handler.Serve(ctx); err != nil {
			r.ene(ctx, err)
		}
	}))
	route := r.Router.Handle(filepath.Join(r.prefix, path), next)
	if utility.NonBlank(method) {
		route.Methods(method)
	}
}

func (r *Router) Handle(path string, handler Handler, middlewares ...Middleware) {
	r.handle("", path, handler, middlewares...)
}

func (r *Router) HandleFunc(path string, handler Handler, middleware ...Middleware) {
	r.handle("", path, handler, middleware...)
}

func (r *Router) HandleMethod(method, path string, handler Handler, middlewares ...Middleware) {
	r.handle(method, path, handler, middlewares...)
}

func (r *Router) HandleMethodFunc(method, path string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(method, path, Handler(handler), middlewares...)
}

func (r *Router) GET(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodGet, path, handler, middlewares...)
}

func (r *Router) POST(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodPost, path, handler, middlewares...)
}

func (r *Router) PUT(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodPut, path, handler, middlewares...)
}

func (r *Router) DELETE(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodDelete, path, handler, middlewares...)
}

func (r *Router) HEAD(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodHead, path, handler, middlewares...)
}

func (r *Router) PATCH(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodPatch, path, handler, middlewares...)
}

func (r *Router) CONNECT(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodConnect, path, handler, middlewares...)
}

func (r *Router) OPTIONS(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodOptions, path, handler, middlewares...)
}

func (r *Router) TRACE(path string, handler Handler, middlewares ...Middleware) {
	r.handle(http.MethodTrace, path, handler, middlewares...)
}

func (r *Router) Group(prefix string, middlewares ...Middleware) *Router {
	route := new(Router)
	route.prefix = filepath.Join(r.prefix, prefix)
	route.Router = r.Router
	route.middlewares = append(route.middlewares, append(r.middlewares, middlewares...)...)
	return route
}
