package middleware

import (
	"fmt"
	"time"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
)

type Logger struct {
}

type logInfo struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Delay    string `json:"delay"`
	Status   int    `json:"status"`
	TraceId  string `json:"trace_id"`
	SpanId   string `json:"span_id"`
	Params   string `json:"params"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

func (l *Logger) Handle(ctx *context.Context) {
	start := time.Now()
	defer func() {
		end := time.Now()
		if !debug.FormatIsJson() {
			ctx.Log.Info(
				"%s %s %.3fms %d %s %s %s", ctx.Request.Method, ctx.Request.URL.Path, float64(end.Sub(start).Microseconds())*0.001, ctx.GetStatus(), ctx.Params.String(), string(ctx.RawContent), string(ctx.RespData),
			)
			return
		}
		ctx.Log.Json(logInfo{
			Method: ctx.Request.Method, Path: ctx.Request.URL.Path, Delay: fmt.Sprintf("%.3fms", float64(end.Sub(start).Microseconds())*0.001), Status: ctx.GetStatus(),
			TraceId: ctx.TraceId(), SpanId: ctx.SpandId(), Params: ctx.Params.String(), Request: string(ctx.RawContent), Response: string(ctx.RespData),
		})
	}()

	ctx.Next()
}
