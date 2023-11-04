package pprof

import (
	"net/http/pprof"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
)

// /debug/pprof/cmdline
type CmdLine struct {
	*controller.Base
}

func NewCmdLine() *CmdLine {
	return &CmdLine{Base: controller.NewBase("default")}
}

func (i *CmdLine) Action(ctx *context.Context) error {
	pprof.Cmdline(ctx.Writer(), ctx.Request)
	return nil
}
