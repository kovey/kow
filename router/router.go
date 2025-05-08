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
