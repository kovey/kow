package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxLen(t *testing.T) {
	m := NewMaxLen()
	ok, err := m.Valid("test", "test", 5)
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = m.Valid("test", "test", 4)
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = m.Valid("test", "test", 3)
	assert.False(t, ok)
	assert.NotNil(t, err)
}
