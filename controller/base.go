package controller

import (
	"github.com/kovey/discovery/grpc"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kow/view"
)

type Base struct {
	v        view.ViewInterface
	services []krpc.ServiceName
	group    string
}

func NewBase(group string, services ...krpc.ServiceName) *Base {
	return &Base{services: services, group: group}
}

func NewBaseBy(tplPath, group string, services ...krpc.ServiceName) *Base {
	b := &Base{v: view.NewDefault(nil), services: services, group: group}
	if err := b.v.Load(tplPath); err != nil {
		panic(err)
	}

	return b
}

func (b *Base) Group() string {
	if b.group == grpc.Str_Empty {
		return grpc.Default
	}

	return b.group
}

func (b *Base) Services() []krpc.ServiceName {
	return b.services
}

func (b *Base) View() view.ViewInterface {
	return b.v
}
