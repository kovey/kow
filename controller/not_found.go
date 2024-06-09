package controller

import (
	"net/http"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
)

type NotFound struct {
	*Base
}

func NewNotFound() *NotFound {
	return &NotFound{Base: NewBase("")}
}

func (n *NotFound) Action(ctx *context.Context) error {
	debug.Dbug("not found: %s", ctx.Request.URL.Path)
	ctx.Status(http.StatusNotFound)
	http.NotFound(ctx.Writer(), ctx.Request)
	return nil
}
