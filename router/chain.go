package router

import (
	"net/http"

	"github.com/kovey/kow/context"
)

type Chain struct {
	Middlewares []context.MiddlewareInterface
	Action      context.ActionInterface
	isFile      bool
	fileServer  http.Handler
}

func (c *Chain) handle(ct *context.Context) {
	if c.isFile {
		c.file(ct)
		return
	}

	ct.Middlerware(c.Middlewares...)
	ct.MiddlerwareStart()
}

func (c *Chain) file(ct *context.Context) {
	ct.Request.URL.Path = ct.Params.GetString("filepath")
	c.fileServer.ServeHTTP(ct.Writer(), ct.Request)
}
