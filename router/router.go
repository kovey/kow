package router

import (
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/validator/rule"
)

type RouterInterface interface {
	Method() string
	Path() string
	Chain() *Chain
	Middleware(...context.MiddlewareInterface) RouterInterface
	Rule(key string, rules ...string) RouterInterface
	Data(data rule.ParamInterface) RouterInterface
}

type New func(m string, p string, a context.ActionInterface) RouterInterface

type NewWith func(m string, p string, a context.Handle) RouterInterface

var NewRouter New = NewDefault
var NewRouterWith NewWith = NewDefaultWith
