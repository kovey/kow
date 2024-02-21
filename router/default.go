package router

import (
	"github.com/kovey/kow/context"
)

type Default struct {
	m     string
	p     string
	chain *Chain
}

func NewDefault(m string, p string, a context.ActionInterface) *Default {
	return &Default{m: m, p: p, chain: &Chain{Action: a}}
}

func (d *Default) Middleware(middlewares ...context.MiddlewareInterface) {
	d.chain.Middlewares = append(d.chain.Middlewares, middlewares...)
}

func (d *Default) Method() string {
	return d.m
}

func (d *Default) Path() string {
	return d.p
}

func (d *Default) Chain() *Chain {
	return d.chain
}
