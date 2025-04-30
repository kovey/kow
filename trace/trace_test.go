package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	data := int64(0)
	en := Encode(data)
	de := Decode(en)
	assert.Equal(t, "A-A", string(en))
	assert.Equal(t, int64(0), de)
	assert.True(t, len(TraceId(10000000000000)) > 0)
}
