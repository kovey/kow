package form

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type test_ref struct {
	Name     string    `form:"name"`
	Id       int64     `form:"id"`
	Sex      int8      `form:"sex"`
	Age      int16     `form:"age"`
	Chips    int32     `form:"chips"`
	TdCount  int       `form:"td_count"`
	Uid      uint64    `form:"uid"`
	USex     uint8     `form:"usex"`
	UAge     uint16    `form:"uage"`
	UChips   uint32    `form:"uchips"`
	UTdCount uint      `form:"utd_count"`
	Bool     bool      `form:"bool"`
	Float32  float32   `form:"float32"`
	Float64  float64   `form:"float64"`
	Date     time.Time `form:"date"`
}

type test_ref_arr struct {
	Name     []string    `form:"name"`
	Id       []int64     `form:"id"`
	Sex      []int8      `form:"sex"`
	Age      []int16     `form:"age"`
	Chips    []int32     `form:"chips"`
	TdCount  []int       `form:"td_count"`
	Uid      []uint64    `form:"uid"`
	USex     []uint8     `form:"usex"`
	UAge     []uint16    `form:"uage"`
	UChips   []uint32    `form:"uchips"`
	UTdCount []uint      `form:"utd_count"`
	Bool     []bool      `form:"bool"`
	Float32  []float32   `form:"float32"`
	Float64  []float64   `form:"float64"`
	Date     []time.Time `form:"date"`
}

func _getTestData() map[string][]string {
	data := make(map[string][]string)
	data["name"] = []string{"kovey"}               // string
	data["id"] = []string{"1"}                     // int64
	data["sex"] = []string{"1"}                    // int8
	data["age"] = []string{"18"}                   // int16
	data["chips"] = []string{"100"}                // int32
	data["td_count"] = []string{"10000"}           // int
	data["uid"] = []string{"1"}                    // uint64
	data["usex"] = []string{"1"}                   // uint8
	data["uage"] = []string{"18"}                  // uint16
	data["uchips"] = []string{"100"}               // uint32
	data["utd_count"] = []string{"10000"}          // uint
	data["bool"] = []string{"true"}                // bool
	data["float32"] = []string{"10000.11"}         // float32
	data["float64"] = []string{"10000.11"}         // float64
	data["date"] = []string{"2025-04-11 11:11:11"} // float64

	return data
}

func TestUnmarshal(t *testing.T) {
	data := _getTestData()
	var tData test_ref
	ref := NewReflect(&tData)
	err := ref.Parse(data, &tData)
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

func TestUnmarshalArr(t *testing.T) {
	data := _getTestData()
	var tData test_ref_arr
	ref := NewReflect(&tData)
	err := ref.Parse(data, &tData)
	assert.Nil(t, err)
	assert.Equal(t, "kovey", tData.Name[0])
	assert.Equal(t, int64(1), tData.Id[0])
	assert.Equal(t, int8(1), tData.Sex[0])
	assert.Equal(t, int16(18), tData.Age[0])
	assert.Equal(t, int32(100), tData.Chips[0])
	assert.Equal(t, int(10000), tData.TdCount[0])
	assert.Equal(t, uint64(1), tData.Uid[0])
	assert.Equal(t, uint8(1), tData.USex[0])
	assert.Equal(t, uint16(18), tData.UAge[0])
	assert.Equal(t, uint32(100), tData.UChips[0])
	assert.Equal(t, uint(10000), tData.UTdCount[0])
	assert.True(t, tData.Bool[0])
	assert.Equal(t, float32(10000.11), tData.Float32[0])
	assert.Equal(t, float64(10000.11), tData.Float64[0])
	assert.Equal(t, "2025-04-11 11:11:11", tData.Date[0].Format(time.DateTime))
}

func TestMarshalArr(t *testing.T) {
	data := _getTestData()
	var tData test_ref_arr
	ref := NewReflect(&tData)
	err := ref.Parse(data, &tData)
	assert.Nil(t, err)

	buf, err := Marshal(tData)
	assert.Nil(t, err)
	assert.Equal(t, "name[]=kovey&id[]=1&sex[]=1&age[]=18&chips[]=100&td_count[]=10000&uid[]=1&usex[]=1&uage[]=18&uchips[]=100&utd_count[]=10000&bool[]=true&float32[]=10000.11035&float64[]=10000.11&date[]=2025-04-11 11:11:11", string(buf))
}

func TestMarshal(t *testing.T) {
	data := _getTestData()
	var tData test_ref
	ref := NewReflect(&tData)
	err := ref.Parse(data, &tData)
	assert.Nil(t, err)

	buf, err := Marshal(tData)
	assert.Nil(t, err)
	assert.Equal(t, "name=kovey&id=1&sex=1&age=18&chips=100&td_count=10000&uid=1&usex=1&uage=18&uchips=100&utd_count=10000&bool=true&float32=10000.11035&float64=10000.11&date=2025-04-11 11:11:11", string(buf))
}
