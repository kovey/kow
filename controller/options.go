package controller

import (
	"net/http"

	"github.com/kovey/kow/context"
)

type Options struct {
	*Base
}

func NewOptions() *Options {
	return &Options{Base: NewBase("")}
}

func (o *Options) Action(ctx *context.Context) error {
	return ctx.Data(http.StatusAccepted, context.Content_Type_Html, nil)
}
