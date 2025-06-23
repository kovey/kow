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

type Handle func(ctx *Context) error

type Action struct {
	action   ActionInterface
	handle   Handle
	services []krpc.ServiceName
	group    string
}

func (a *Action) WithAction(action ActionInterface) *Action {
	a.action = action
	a.services = action.Services()
	a.group = action.Group()
	return a
}

func (a *Action) WithServices(services ...krpc.ServiceName) *Action {
	a.services = append(a.services, services...)
	return a
}

func (a *Action) WithGroup(group string) *Action {
	a.group = group
	return a
}

func (a *Action) WithHandle(handle Handle) *Action {
	a.handle = handle
	return a
}

func (a *Action) Action(ctx *Context) error {
	if a.action != nil {
		return a.action.Action(ctx)
	}

	return a.handle(ctx)
}

func (a *Action) Services() []krpc.ServiceName {
	return a.services
}

func (a *Action) Group() string {
	return a.group
}
