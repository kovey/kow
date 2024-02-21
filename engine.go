package kow

import (
	cc "context"
	"net/http"
	"time"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/middleware"
	"github.com/kovey/kow/router"
)

type Engine struct {
	routers    *router.Routers
	maxRunTime time.Duration
	serv       *http.Server
}

func NewEngine() *Engine {
	return &Engine{routers: router.NewRouters()}
}

func NewDefault() *Engine {
	e := NewEngine()
	e.Middleware(&middleware.Logger{}, &middleware.Recovery{})
	return e
}

func (e *Engine) SetMaxRunTime(max time.Duration) {
	e.maxRunTime = max
}

func (e *Engine) GET(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodGet, path, ac)
}

func (e *Engine) POST(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodPost, path, ac)
}

func (e *Engine) PUT(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodPut, path, ac)
}

func (e *Engine) PATCH(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodPatch, path, ac)
}

func (e *Engine) HEAD(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodHead, path, ac)
}

func (e *Engine) DELETE(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodDelete, path, ac)
}

func (e *Engine) CONNECT(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodConnect, path, ac)
}

func (e *Engine) OPTIONS(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodOptions, path, ac)
}

func (e *Engine) TRACE(path string, ac context.ActionInterface) router.RouterInterface {
	return e.DefRouter(http.MethodTrace, path, ac)
}

func (e *Engine) DefRouter(method string, path string, ac context.ActionInterface) router.RouterInterface {
	ro := router.NewDefault(method, path, ac)
	e.routers.Add(ro)
	return ro
}

func (e *Engine) Router(ro router.RouterInterface) {
	e.routers.Add(ro)
}

func (e *Engine) Middleware(middlewars ...context.MiddlewareInterface) {
	e.routers.Middlerware(middlewars...)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parent, cancel := cc.WithTimeout(r.Context(), e.maxRunTime)
	ctx := context.NewContext(parent, w, r)
	defer ctx.Drop()
	defer cancel()

	e.routers.HandleHTTP(ctx)
}

func (e *Engine) Run(addr string) error {
	e.serv = &http.Server{Addr: addr, Handler: e}
	return e.serv.ListenAndServe()
}

func (e *Engine) Shutdown() error {
	if e.serv == nil {
		return nil
	}

	return e.serv.Shutdown(cc.Background())
}
