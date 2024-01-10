package resolver

import (
	"github.com/kovey/discovery/etcd"
	"github.com/kovey/discovery/resolver"
)

func Register(conf etcd.Config) error {
	resolver.Init(conf)
	return resolver.Register()
}

func Shutdown() {
	resolver.Shutdown()
}
