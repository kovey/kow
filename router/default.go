package router

import (
	"github.com/kovey/kow/context"
)

type Default struct {
	M     string
	P     string
	chain *Chain
}

func NewDefault(m string, p string, a context.ActionInterface) *Default {
	return &Default{M: m, P: p, chain: &Chain{Action: a}}
}

func (d *Default) Middleware(middlewares ...context.MiddlewareInterface) {
	d.chain.Middlewares = append(d.chain.Middlewares, middlewares...)
}

func (d *Default) Method() string {
	return d.M
}

func (d *Default) Path() string {
	return d.P
}

func (d *Default) Chain() *Chain {
	return d.chain
}
