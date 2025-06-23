package kow

import (
	cc "context"
	"net/http"
	"time"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
	"github.com/kovey/kow/middleware"
	"github.com/kovey/kow/router"
	"github.com/kovey/pool"
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
	e.Middleware(&middleware.Logger{}, &middleware.Recovery{}, middleware.NewParseRequestData(), middleware.NewValidator())
	e.routers.NotFound = router.NewChain(controller.NewNotFound())
	e.routers.GlobalOPTIONS = router.NewChain(controller.NewOptions())
	e.routers.HandleOPTIONS = true
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
	ro := router.NewRouter(method, path, ac)
	e.routers.Add(ro)
	return ro
}

func (e *Engine) RouterWith(method string, path string, handle context.Handle) router.RouterInterface {
	ro := router.NewRouterWith(method, path, handle)
	e.routers.Add(ro)
	return ro
}

func (e *Engine) Group(path string) *router.Group {
	return e.routers.Group(path)
}

func (e *Engine) Router(ro router.RouterInterface) {
	e.routers.Add(ro)
}

func (e *Engine) Middleware(middlewars ...context.MiddlewareInterface) {
	e.routers.Middlerware(middlewars...)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	maxRunTime := e.maxRunTime
	if maxRunTime <= 0 {
		maxRunTime = 60 * time.Second
	}
	parent, cancel := cc.WithTimeout(r.Context(), maxRunTime)
	ctx := pool.NewContext(parent)
	defer ctx.Drop()
	defer cancel()

	e.routers.HandleHTTP(context.NewContext(ctx, w, r))
}

func (e *Engine) Run(addr string) error {
	e.serv = &http.Server{Addr: addr, Handler: e}
	err := e.serv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (e *Engine) Shutdown() error {
	if e.serv == nil {
		return nil
	}

	return e.serv.Shutdown(cc.Background())
}
