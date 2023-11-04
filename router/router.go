package router

import (
	"github.com/kovey/kow/context"
)

type RouterInterface interface {
	Method() string
	Path() string
	Chain() *Chain
	Middleware(...context.MiddlewareInterface)
}
