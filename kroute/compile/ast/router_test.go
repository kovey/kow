package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	r := &Router{}
	content := []byte(`("PUT", "/path", newAction(), &ReqData{})`)
	for _, b := range content {
		assert.Nil(t, r.Parse(b))
	}

	assert.Equal(t, "PUT", r.Method.Value)
	assert.Equal(t, byte('"'), r.Method.Begin)
	assert.Equal(t, byte('"'), r.Method.End)
	assert.Equal(t, byte('('), r.Begin)
	assert.Equal(t, byte(')'), r.End)
	assert.Equal(t, "/path", r.Path.Value)
	assert.Equal(t, byte('"'), r.Path.Begin)
	assert.Equal(t, byte('"'), r.Path.End)
	assert.NotNil(t, r.Constructor)
	assert.Equal(t, byte('('), r.Constructor.Left)
	assert.Equal(t, "newAction", r.Constructor.Name)
	assert.Equal(t, byte(')'), r.Constructor.Right)
	assert.NotNil(t, r.ReqData)
	assert.Equal(t, byte('&'), r.ReqData.Star)
	assert.Equal(t, byte('{'), r.ReqData.Left)
	assert.Equal(t, "ReqData", r.ReqData.Name)
	assert.Equal(t, byte('}'), r.ReqData.Right)
}
