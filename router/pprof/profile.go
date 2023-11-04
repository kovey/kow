package pprof

import (
	"net/http/pprof"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
)

// /debug/pprof/profile
type Profile struct {
	*controller.Base
}

func NewProfile() *Profile {
	return &Profile{Base: controller.NewBase("default")}
}

func (i *Profile) Action(ctx *context.Context) error {
	pprof.Profile(ctx.Writer(), ctx.Request)
	return nil
}
