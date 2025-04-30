package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	gt := NewLen()
	ok, err := gt.Valid("test", 1)
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = gt.Valid("test", 1, 1)
	assert.False(t, ok)
	assert.NotNil(t, err)

	ok, err = gt.Valid("test", "1", 2)
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = gt.Valid("test", "2", 1)
	assert.True(t, ok)
	assert.Nil(t, err)
}
