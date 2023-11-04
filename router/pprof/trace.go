package pprof

import (
	"net/http/pprof"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
)

// debug/pprof/trace
type Trace struct {
	*controller.Base
}

func NewTrace() *Trace {
	return &Trace{Base: controller.NewBase("default")}
}

func (i *Trace) Action(ctx *context.Context) error {
	pprof.Trace(ctx.Writer(), ctx.Request)
	return nil
}
