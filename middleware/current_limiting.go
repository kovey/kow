package middleware

import (
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/funnel"
)

type CurrentLimiting struct {
}

func (c *CurrentLimiting) Handle(ctx *context.Context) {
	funnel.Get()
	ctx.Next()
}
