package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRouter(t *testing.T) {
	d := NewDefault("POST", "/", newTestAction())
	d.Rule("email", "email").Rule("age", "gt:int:0", "le:int:120")
	d.Middleware(&test_middle{})
	d.Data(&req_data{})
	assert.Equal(t, "POST", d.Method())
	assert.Equal(t, "/", d.Path())
	assert.Equal(t, 1, len(d.Chain().rules.Get("email")))
	assert.Equal(t, 2, len(d.Chain().rules.Get("age")))
	assert.Equal(t, 1, len(d.Chain().Middlewares))
	assert.NotNil(t, d.chain.param)
}
