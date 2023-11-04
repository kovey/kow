package main

import (
	"net/http"
	"time"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/etcd"
	"github.com/kovey/kow"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
	"github.com/kovey/kow/resolver"
	"github.com/kovey/kow/tests/protocol"
	"github.com/kovey/kow/tests/service"
	"github.com/kovey/kow/view"
)

type Middle struct {
}

func (m *Middle) Handle(ctx *context.Context) {
	if ctx.Params.GetString("name") != "1" {
		debug.Erro("name[%s] is not 1", ctx.Params.GetString("name"))
		ctx.Json(http.StatusOK, ctx.Params)
		return
	}

	ctx.Next()
}

type Index struct {
	*controller.Base
}

func NewIndex() *Index {
	return &Index{Base: controller.NewBase("test", service.Service_Test)}
}

func (i *Index) Action(ctx *context.Context) error {
	c := protocol.NewHelloClient(ctx.Rpcs.Get(service.Service_Test, "test"))
	res, err := c.Error(ctx, &protocol.HelloInfo{Name: "name"})
	if err != nil {
		ctx.Json(http.StatusOK, map[string]string{"msg": err.Error()})
		return err
	}
	ctx.Json(http.StatusOK, map[string]string{"kovey": res.Name})
	return nil
}

func (i *Index) View() view.ViewInterface {
	return nil
}

type IndexTest struct {
	*controller.Base
}

func NewIndexTest() *IndexTest {
	return &IndexTest{Base: controller.NewBase("test", service.Service_Test)}
}

func (i *IndexTest) Action(ctx *context.Context) error {
	c := protocol.NewHelloClient(ctx.Rpcs.Get(service.Service_Test, "test"))
	res, err := c.Say(ctx, &protocol.HelloInfo{Name: ctx.Params.GetString("name")})
	if err != nil {
		ctx.Json(http.StatusOK, map[string]string{"msg": err.Error()})
	}

	ctx.Json(http.StatusOK, map[string]string{"kovey": res.Name})
	return nil
}

func (i *IndexTest) View() view.ViewInterface {
	return nil
}
func main() {
	e := kow.NewDefault()
	e.SetMaxRunTime(5 * time.Second)
	e.GET("/index", NewIndex())
	e.GET("/test/:name/:user", NewIndexTest()).Middleware(&Middle{})

	if err := resolver.Register(etcd.Config{Endpoints: []string{"47.108.158.214:9001"}, DialTimeout: 30}); err != nil {
		panic(err)
	}

	if err := e.Run("0.0.0.0:8081"); err != nil {
		panic(err)
	}
}
