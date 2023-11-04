package pprof

import (
	"net/http/pprof"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
)

// /debug/pprof/symbol
type Symbol struct {
	*controller.Base
}

func NewSymbol() *Symbol {
	return &Symbol{Base: controller.NewBase("default")}
}

func (i *Symbol) Action(ctx *context.Context) error {
	pprof.Symbol(ctx.Writer(), ctx.Request)
	return nil
}
