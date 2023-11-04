package pprof

import (
	"net/http/pprof"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
)

// /debug/pprof/
type Index struct {
	*controller.Base
}

func NewIndex() *Index {
	return &Index{Base: controller.NewBase("default")}
}

func (i *Index) Action(ctx *context.Context) error {
	pprof.Index(ctx.Writer(), ctx.Request)
	return nil
}
