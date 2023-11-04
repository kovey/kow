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
		debug.Info("%s %s %.2fms %d", ctx.Request.Method, ctx.Request.URL.Path, float64(end.Sub(start).Microseconds())*0.0001, ctx.GetStatus())
	}()

	ctx.Next()
}
