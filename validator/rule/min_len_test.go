package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinLen(t *testing.T) {
	m := NewMinLen()
	ok, err := m.Valid("test", "test", 5)
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = m.Valid("test", "test", 4)
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = m.Valid("test", "test", 3)
	assert.True(t, ok)
	assert.Nil(t, err)
}
