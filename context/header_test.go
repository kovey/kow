package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	assert.True(t, IsJson("application/json"))
	assert.False(t, IsJson("application/xml"))
	assert.True(t, IsBinary("application/octet-stream"))
	assert.False(t, IsBinary("application/xml"))
	assert.True(t, IsForm("application/x-www-form-urlencoded"))
	assert.False(t, IsForm("application/xml"))
	assert.True(t, IsXml("text/xml"))
	assert.False(t, IsXml("application/xml"))
}
