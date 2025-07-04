package kow

import (
	"net/http"
	"time"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
	"github.com/kovey/kow/middleware"
	"github.com/kovey/kow/router"
	"github.com/kovey/kow/serv"
)

const (
	APP_TIME_ZONE  = "APP_TIME_ZONE"
	APP_ETCD_OPEN  = "APP_ETCD_OPEN"
	APP_PPROF_OPEN = "APP_PPROF_OPEN"
	SERV_HOST      = "SERV_HOST"
	SERV_PORT      = "SERV_PORT"
	ETCD_TIMEOUT   = "ETCD_TIMEOUT"
	ETCD_USERNAME  = "ETCD_USERNAME"
	ETCD_PASSWORD  = "ETCD_PASSWORD"
	ETCD_NAMESPACE = "ETCD_NAMESPACE"
	ETCD_ENDPOINTS = "ETCD_ENDPOINTS"
)

func SetMaxRunTime(max time.Duration) {
	engine.SetMaxRunTime(max)
}

func GET(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodGet, path, ac)
}

func POST(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodPost, path, ac)
}

func PUT(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodPut, path, ac)
}

func PATCH(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodPatch, path, ac)
}

func HEAD(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodHead, path, ac)
}

func DELETE(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodDelete, path, ac)
}

func CONNECT(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodConnect, path, ac)
}

func OPTIONS(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodOptions, path, ac)
}

func TRACE(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodTrace, path, ac)
}

func Group(path string) *router.Group {
	return engine.Group(path)
}

func Router(method, path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(method, path, ac)
}

func RouterWith(method, path string, handle context.Handle) router.RouterInterface {
	return engine.RouterWith(method, path, handle)
}

func Run(e serv.EventInterface) {
	name := "kow"
	if e != nil {
		name = e.AppName()
	}

	cli := app.NewApp(name)
	serv := newServer(e)
	cli.SetServ(serv)
	if err := cli.Run(); err != nil {
		debug.Erro(err.Error())
	}
}

func OpenCors(headers ...string) {
	engine.routers.HandleOPTIONS = true
	engine.routers.GlobalOPTIONS = router.NewChain(controller.NewOptions())
	engine.Middleware(&middleware.OpenCors{Headers: headers})
}

func Middleware(m ...context.MiddlewareInterface) {
	engine.Middleware(m...)
}

func SetGlobalOPTIONS(act context.ActionInterface) {
	if engine.routers.GlobalOPTIONS == nil {
		engine.routers.GlobalOPTIONS = &router.Chain{}
	}

	ac := &context.Action{}
	engine.routers.GlobalOPTIONS.Action = ac.WithAction(act)
}

func SetNotFound(act context.ActionInterface) {
	if engine.routers.NotFound == nil {
		engine.routers.NotFound = &router.Chain{}
	}

	ac := &context.Action{}
	engine.routers.NotFound.Action = ac.WithAction(act)
}
