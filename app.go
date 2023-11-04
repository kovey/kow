package kow

import (
	"net/http"
	"time"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/middleware"
	"github.com/kovey/kow/router"
	"github.com/kovey/kow/serv"
)

func SetMaxRunTime(max time.Duration) {
	engine.maxRunTime = max
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

func Patch(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodPatch, path, ac)
}

func HEAD(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodHead, path, ac)
}

func DELETE(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodDelete, path, ac)
}

func Connect(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodConnect, path, ac)
}

func OPTIONS(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodOptions, path, ac)
}

func TRACE(path string, ac context.ActionInterface) router.RouterInterface {
	return engine.DefRouter(http.MethodTrace, path, ac)
}

func Run(name string, e serv.EventInterface) {
	cli := app.NewApp(name)
	serv := newServer(e)
	cli.SetDebugLevel(debug.Debug_Info)
	cli.SetServ(serv)
	if err := cli.Run(); err != nil {
		debug.Erro(err.Error())
	}
}

func OpenCors(headers ...string) {
	engine.routers.HandleOPTIONS = true
	engine.routers.GlobalOPTIONS = &router.Chain{}
	engine.Middleware(&middleware.OpenCors{Headers: headers})
}

func Middleware(m ...context.MiddlewareInterface) {
	engine.Middleware(m...)
}
