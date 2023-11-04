package middleware

import (
	"strings"

	"github.com/kovey/kow/context"
)

type OpenCors struct {
	Headers []string
}

func (o *OpenCors) Handle(ctx *context.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", o.get())
	ctx.Next()
}

func (o *OpenCors) def() []string {
	return []string{"content-type"}
}

func (o *OpenCors) get() string {
	return strings.Join(append(o.def(), o.Headers...), ",")
}
