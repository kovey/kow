package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	data := int64(0)
	en := Encode(data)
	de := Decode(en)
	assert.Equal(t, "AAAAAAA-AAAAAAA", string(en))
	assert.Equal(t, int64(0), de)
	assert.Panics(t, func() {
		Decode([]byte("AAAA-AAAAAAA"))
	})
	assert.Panics(t, func() {
		Decode([]byte("AAAAOIB-AAAAAAA"))
	})
	traceId := TraceId(10000000000000)
	assert.True(t, len(traceId) > 0)
	t.Logf("traceId: %s", traceId)
}
