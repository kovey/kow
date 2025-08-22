package middleware

import (
	"net/http"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/funnel"
)

type CurrentLimiting struct {
}

func (c *CurrentLimiting) Handle(ctx *context.Context) {
	f := funnel.Get()
	if f == 0 {
		ctx.Html(http.StatusServiceUnavailable, nil)
		return
	}

	ctx.Next()
}
