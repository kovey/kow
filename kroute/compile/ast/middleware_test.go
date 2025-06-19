package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	m := &Middleware{}
	content := []byte(`(newAuth(), middlewares.NewLog())`)
	for _, b := range content {
		assert.Nil(t, m.Parse(b))
	}

	assert.Len(t, m.Args, 2)
	assert.Equal(t, byte('('), m.Begin)
	assert.Equal(t, byte(')'), m.End)
	assert.Equal(t, byte('('), m.Args[0].Left)
	assert.Equal(t, "newAuth", m.Args[0].Name)
	assert.Equal(t, byte(')'), m.Args[0].Right)
	assert.Equal(t, byte('('), m.Args[1].Left)
	assert.Equal(t, "middlewares.NewLog", m.Args[1].Name)
	assert.Equal(t, byte(')'), m.Args[1].Right)
}

func TestMiddlewareError(t *testing.T) {
	m := &Middleware{}
	content := []byte(`(newAuth(,), middlewares.NewLog())`)
	var err error
	for _, b := range content {
		err = m.Parse(b)
		if err != nil {
			break
		}
	}

	assert.NotNil(t, err)
	assert.Equal(t, "middleware unexpect ,", err.Error())
}

func TestMiddlewareErrorFormat(t *testing.T) {
	m := &Middleware{}
	content := []byte(`(newAuth(), , middlewares.NewLog())`)
	var err error
	for _, b := range content {
		err = m.Parse(b)
		if err != nil {
			break
		}
	}

	assert.NotNil(t, err)
	assert.Equal(t, "middleware unexpect ,", err.Error())
}
