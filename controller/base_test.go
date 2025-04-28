package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	b := NewBase("group")
	assert.Equal(t, "group", b.Group())
	assert.Nil(t, b.Services())
	assert.Nil(t, b.View())
}

func TestBaseBy(t *testing.T) {
	b := NewBaseBy("./test_tpl.html", "")
	assert.Equal(t, "default", b.Group())
	assert.Nil(t, b.Services())
	assert.NotNil(t, b.View())
}
