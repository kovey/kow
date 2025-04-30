package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChinese(t *testing.T) {
	c := NewChinese()
	ok, err := c.Valid("key", "中文")
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = c.Valid("key", "english")
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = c.Valid("key", 1)
	assert.False(t, ok)
}
