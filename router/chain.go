package router

import (
	"net/http"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/validator"
	"github.com/kovey/kow/validator/rule"
)

type Chain struct {
	Middlewares []context.MiddlewareInterface
	Action      *context.Action
	isFile      bool
	fileServer  http.Handler
	rules       *validator.ParamRules
	param       rule.ParamInterface
}

func NewChain(act context.ActionInterface) *Chain {
	ac := &context.Action{}
	return &Chain{Action: ac.WithAction(act)}
}

func (c *Chain) SetFileServer(fileServer http.Handler) {
	c.isFile = true
	c.fileServer = fileServer
}

func (c *Chain) handle(ct *context.Context) {
	if c.isFile {
		c.file(ct)
		return
	}

	ct.Rules = c.rules
	if c.param != nil {
		ct.ReqData = c.param.Clone()
	}
	ct.SetAction(c.Action)
	ct.Middleware(c.Middlewares...)
	ct.MiddlerwareStart()
}

func (c *Chain) file(ct *context.Context) {
	ct.Request.URL.Path = ct.Params.GetString("filepath")
	c.fileServer.ServeHTTP(ct.Writer(), ct.Request)
}
