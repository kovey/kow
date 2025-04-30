package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomain(t *testing.T) {
	c := NewDomain()
	ok, err := c.Valid("key", "www.baidu.com")
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = c.Valid("key", "https://www.baidu.com")
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = c.Valid("key", 1)
	assert.False(t, ok)
}
