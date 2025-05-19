package middleware

import (
	"time"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
)

type Logger struct {
}

func (l *Logger) Handle(ctx *context.Context) {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.Info("%s %s %.3fms %d %s %s", ctx.Request.Method, ctx.Request.URL.Path, float64(end.Sub(start).Microseconds())*0.001, ctx.GetStatus(), ctx.TraceId(), string(ctx.RawContent))
	}()

	ctx.Next()
}
