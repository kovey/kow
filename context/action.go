package context

import (
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kow/view"
	"google.golang.org/grpc"
)

type ActionInterface interface {
	Action(c *Context) error
	View() view.ViewInterface
	Services() []krpc.ServiceName
	Group() string
}

type Rpcs map[string]grpc.ClientConnInterface

func (r Rpcs) Get(serviceName krpc.ServiceName, group string) grpc.ClientConnInterface {
	return r[serviceName.Group(group)]
}

func (r Rpcs) Default(serviceName krpc.ServiceName) grpc.ClientConnInterface {
	return r[serviceName.Default()]
}

func (r Rpcs) Add(serviceName krpc.ServiceName, group string, conn grpc.ClientConnInterface) {
	r[serviceName.Group(group)] = conn
}

func (r Rpcs) AddDefualt(serviceName krpc.ServiceName, conn grpc.ClientConnInterface) {
	r[serviceName.Default()] = conn
}
