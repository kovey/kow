package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamRules(t *testing.T) {
	r := NewParamRules()
	r.Add("nickname", "eq:string:kovey")
	r.Add("name", "maxlen:int:20", "minlen:int:6")
	r.Add("status", "eq:int8:20")
	r.Add("sex", "eq:int16:20")
	r.Add("age", "eq:int32:20")
	r.Add("balance", "eq:int64:20")
	r.Add("name1", "maxlen:uint:20", "minlen:int:6")
	r.Add("status1", "eq:uint8:20")
	r.Add("sex1", "eq:uint16:20")
	r.Add("age1", "eq:uint32:20")
	r.Add("balance1", "eq:uint64:20")
	r.Add("has_other", "eq:bool:true")
	r.Add("money", "eq:float32:10.12")
	r.Add("money1", "eq:float64:100.12")
	r.Add("regx_test", `regx:string:[\\u4e00-\\u9fa5]`)

	assert.False(t, r.Add("name", "eq:string:kovey"))
	assert.False(t, r.Add("name", "eq:string"))
	assert.True(t, len(r.rules) == 15)
	assert.True(t, len(r.Get("name")) == 2)
	assert.True(t, len(r.Get("status")) == 1)
}
