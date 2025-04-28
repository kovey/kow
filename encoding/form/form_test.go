package form

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormUnmarshal(t *testing.T) {
	data := _getTestData()
	var tData test_ref
	err := Unmarshal(data, &tData)
	assert.Nil(t, err)
	assert.Equal(t, "kovey", tData.Name)
	assert.Equal(t, int64(1), tData.Id)
	assert.Equal(t, int8(1), tData.Sex)
	assert.Equal(t, int16(18), tData.Age)
	assert.Equal(t, int32(100), tData.Chips)
	assert.Equal(t, int(10000), tData.TdCount)
	assert.Equal(t, uint64(1), tData.Uid)
	assert.Equal(t, uint8(1), tData.USex)
	assert.Equal(t, uint16(18), tData.UAge)
	assert.Equal(t, uint32(100), tData.UChips)
	assert.Equal(t, uint(10000), tData.UTdCount)
	assert.True(t, tData.Bool)
	assert.Equal(t, float32(10000.11), tData.Float32)
	assert.Equal(t, float64(10000.11), tData.Float64)
	assert.Equal(t, "2025-04-11 11:11:11", tData.Date.Format(time.DateTime))
}
