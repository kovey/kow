package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeMap(t *testing.T) {
	sm := NewSafeMap[string, int]()
	assert.Equal(t, 0, sm.Len())
	sm.Add("age", 18)
	assert.Equal(t, 1, sm.Len())
	sm.Add("count", 100)
	assert.Equal(t, 2, sm.Len())
	assert.Equal(t, 18, sm.Get("age"))
	assert.Equal(t, 100, sm.Get("count"))
	sm.Rem("age")
	assert.Equal(t, 0, sm.Get("age"))
	assert.False(t, sm.Exists("age"))
}

func TestSafeMapRange(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Add("age", 18)
	sm.Add("count", 100)

	sm.Range(func(key string, val int) bool {
		switch key {
		case "age":
			assert.Equal(t, 18, val)
			return true
		case "count":
			assert.Equal(t, 100, val)
			return true
		}
		return true
	})
}
