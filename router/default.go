package router

import (
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/validator"
)

type Default struct {
	m     string
	p     string
	chain *Chain
}

func NewDefault(m string, p string, a context.ActionInterface) *Default {
	return &Default{m: m, p: p, chain: &Chain{Action: a, rules: validator.NewParamRules()}}
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

func (d *Default) Rule(key string, rules ...string) RouterInterface {
	d.chain.rules.Add(key, rules...)
	return d
}
