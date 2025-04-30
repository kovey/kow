package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrl(t *testing.T) {
	c := NewUrl()
	ok, err := c.Valid("key", "https://www.baidu.com")
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = c.Valid("key", "kovey.com")
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = c.Valid("key", 1)
	assert.False(t, ok)
}
