package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData(t *testing.T) {
	dt := Params{"username": "kovey", "age": "18", "balance": "10.23", "count": "999", "is_man": "true", MatchedRoutePathParam: "index"}
	dt.MatchedRoutePath()
	assert.Equal(t, "index", dt.MatchedRoutePath())
	assert.Equal(t, "kovey", dt.GetString("username"))
	assert.Equal(t, 18, dt.GetInt("age"))
	assert.Equal(t, 0, dt.GetInt("username"))
	assert.Equal(t, float64(10.23), dt.GetFloat("balance"))
	assert.Equal(t, float64(0), dt.GetFloat("username"))
	assert.Equal(t, int64(999), dt.GetInt64("count"))
	assert.Equal(t, int64(0), dt.GetInt64("username"))
	assert.True(t, dt.GetBool("is_man"))
	assert.False(t, dt.GetBool("username"))
	dt.Reset()
	assert.Equal(t, "", dt.GetString("username"))
}
