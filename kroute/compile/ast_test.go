package compile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAst(t *testing.T) {
	a := &Ast{}
	content := []string{
		`//go:kroute router("PUT", "/path", newAction(), &ReqData{}) rule("name", ["minlen:int:1", "maxlen:int:127"]) middleware(newAuth(), middlewares.NewLog()) import("path1", "path2")`,
		`//go:kroute rule("age", ["minlen:int:1", "maxlen:int:127", "email"])`,
		`//go:kroute middleware(newSign())`,
	}

	assert.Nil(t, a.Parse(content))

	assert.NotNil(t, a.Router)
	assert.NotNil(t, a.Router.Method)
	assert.NotNil(t, a.Router.Path)
	assert.NotNil(t, a.Router.Constructor)
	assert.NotNil(t, a.Router.ReqData)
	assert.Equal(t, "PUT", a.Router.Method.Value)
	assert.Equal(t, "/path", a.Router.Path.Value)
	assert.Equal(t, "newAction", a.Router.Constructor.Name)
	assert.Equal(t, byte('&'), a.Router.ReqData.Star)
	assert.Equal(t, "ReqData", a.Router.ReqData.Name)

	assert.Len(t, a.Imports, 1)
	assert.Len(t, a.Imports[0].Args, 2)
	assert.Equal(t, "path1", a.Imports[0].Args[0].Value)
	assert.Equal(t, "path2", a.Imports[0].Args[1].Value)

	assert.Len(t, a.Middlewares, 2)
	assert.Len(t, a.Middlewares[0].Args, 2)
	assert.Len(t, a.Middlewares[1].Args, 1)
	assert.Len(t, a.Middlewares, 2)
	assert.Equal(t, "newAuth", a.Middlewares[0].Args[0].Name)
	assert.Equal(t, "middlewares.NewLog", a.Middlewares[0].Args[1].Name)
	assert.Equal(t, "newSign", a.Middlewares[1].Args[0].Name)

	assert.Len(t, a.Rules, 2)
	assert.Len(t, a.Rules[0].Args.Args, 2)
	assert.Len(t, a.Rules[1].Args.Args, 3)
	assert.Equal(t, "name", a.Rules[0].Name.Value)
	assert.Equal(t, "minlen:int:1", a.Rules[0].Args.Args[0].Value)
	assert.Equal(t, "maxlen:int:127", a.Rules[0].Args.Args[1].Value)
	assert.Equal(t, "age", a.Rules[1].Name.Value)
	assert.Equal(t, "minlen:int:1", a.Rules[1].Args.Args[0].Value)
	assert.Equal(t, "maxlen:int:127", a.Rules[1].Args.Args[1].Value)
	assert.Equal(t, "email", a.Rules[1].Args.Args[2].Value)
}
