package router

import (
	"fmt"
	"net/http"

	"github.com/kovey/kow/context"
)

type Group struct {
	middlewares []context.MiddlewareInterface
	path        string
	routers     *Routers
}

func NewGroup(path string, routers *Routers) *Group {
	return &Group{path: path, routers: routers}
}

func (g *Group) Middlewares() []context.MiddlewareInterface {
	return g.middlewares
}

func (g *Group) Middleware(middlewares ...context.MiddlewareInterface) *Group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

func (g *Group) Path() string {
	return g.path
}

func (g *Group) GET(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodGet, path, ac)
}

func (g *Group) POST(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodPost, path, ac)
}

func (g *Group) PUT(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodPut, path, ac)
}

func (g *Group) PATCH(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodPatch, path, ac)
}

func (g *Group) HEAD(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodHead, path, ac)
}

func (g *Group) DELETE(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodDelete, path, ac)
}

func (g *Group) CONNECT(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodConnect, path, ac)
}

func (g *Group) OPTIONS(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodOptions, path, ac)
}

func (g *Group) TRACE(path string, ac context.ActionInterface) RouterInterface {
	return g.DefRouter(http.MethodTrace, path, ac)
}

func (g *Group) DefRouter(method string, path string, ac context.ActionInterface) RouterInterface {
	ro := NewDefault(method, fmt.Sprintf("%s/%s", g.path, path), ac)
	ro.Middleware(g.middlewares...)
	g.routers.Add(ro)
	return ro
}
