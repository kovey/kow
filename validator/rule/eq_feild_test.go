package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqFeild(t *testing.T) {
	eq := NewEqFeild()
	ok, err := eq.Valid("test", 1)
	assert.False(t, ok)
	assert.NotNil(t, err)

	ok, err = eq.Valid("test", 1, "name")
	assert.False(t, ok)
	assert.NotNil(t, err)

	ok, err = eq.Valid("test", 1, 1)
	assert.True(t, ok)
	assert.Nil(t, err)
}
