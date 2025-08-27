package middleware

import (
	"net/http"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/funnel"
)

type CurrentLimiting struct {
	name string
}

func NewCurrentLimiting(name string) *CurrentLimiting {
	return &CurrentLimiting{name: name}
}

func (c *CurrentLimiting) Handle(ctx *context.Context) {
	f := funnel.Get(c.name)
	if f == 0 {
		ctx.Html(http.StatusServiceUnavailable, nil)
		return
	}

	ctx.Next()
}
