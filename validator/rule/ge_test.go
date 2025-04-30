package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ge := NewGe()
	ok, err := ge.Valid("test", 1)
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = ge.Valid("test", 1, 1)
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = ge.Valid("test", 2.2, 1.1)
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = ge.Valid("test", 1, 2)
	assert.False(t, ok)
	assert.NotNil(t, err)
}
