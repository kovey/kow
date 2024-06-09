package controller

import (
	"net/http"

	"github.com/kovey/kow/context"
)

type NotFound struct {
	*Base
}

func NewNotFound() *NotFound {
	return &NotFound{Base: NewBase("")}
}

func (n *NotFound) Action(ctx *context.Context) error {
	ctx.Status(http.StatusNotFound)
	http.NotFound(ctx.Writer(), ctx.Request)
	return nil
}
