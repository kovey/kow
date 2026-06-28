package uri

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type uriTestStruct struct {
	ID     int    `uri:"id"`
	Name   string `uri:"name"`
	Email  string `uri:"email"`
	Active bool   `uri:"active"`
}

func TestUnmarshal_Basic(t *testing.T) {
	var s uriTestStruct
	data := map[string]string{
		"id":     "42",
		"name":   "alice",
		"active": "true",
	}

	require.NoError(t, Unmarshal(data, &s))
	assert.Equal(t, 42, s.ID)
	assert.Equal(t, "alice", s.Name)
	assert.Equal(t, true, s.Active)
	assert.Empty(t, s.Email) // not in data
}

func TestUnmarshal_MissingTag(t *testing.T) {
	type testStruct struct {
		NoTag string `json:"name"`
	}
	var s testStruct
	data := map[string]string{"name": "alice"}
	require.NoError(t, Unmarshal(data, &s))
	assert.Empty(t, s.NoTag) // no uri tag
}

func TestUnmarshal_SkipDash(t *testing.T) {
	type testStruct struct {
		Skip string `uri:"-"`
	}
	var s testStruct
	data := map[string]string{"-": "alice"}
	require.NoError(t, Unmarshal(data, &s))
	assert.Empty(t, s.Skip)
}

func TestUnmarshal_NilPointer(t *testing.T) {
	err := Unmarshal(nil, nil)
	assert.Error(t, err)
}

func TestUnmarshal_NotPointer(t *testing.T) {
	var s uriTestStruct
	err := Unmarshal(map[string]string{"id": "1"}, s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not pointer")
}

func TestUnmarshal_NotStruct(t *testing.T) {
	var id int
	err := Unmarshal(map[string]string{"id": "1"}, &id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not struct")
}

func TestUnmarshal_EmptyData(t *testing.T) {
	var s uriTestStruct
	s.ID = 99
	require.NoError(t, Unmarshal(map[string]string{}, &s))
	assert.Equal(t, 99, s.ID) // unchanged
}

func TestUnmarshal_IntTypes(t *testing.T) {
	type testStruct struct {
		I   int   `uri:"i"`
		I8  int8  `uri:"i8"`
		I16 int16 `uri:"i16"`
		I32 int32 `uri:"i32"`
		I64 int64 `uri:"i64"`
	}
	var s testStruct
	data := map[string]string{
		"i":   "-1",
		"i8":  "127",
		"i16": "32767",
		"i32": "2147483647",
		"i64": "9223372036854775807",
	}
	require.NoError(t, Unmarshal(data, &s))
	assert.Equal(t, -1, s.I)
	assert.Equal(t, int8(127), s.I8)
	assert.Equal(t, int16(32767), s.I16)
	assert.Equal(t, int32(2147483647), s.I32)
	assert.Equal(t, int64(9223372036854775807), s.I64)
}

func TestUnmarshal_UintTypes(t *testing.T) {
	type testStruct struct {
		U   uint   `uri:"u"`
		U8  uint8  `uri:"u8"`
		U16 uint16 `uri:"u16"`
		U32 uint32 `uri:"u32"`
		U64 uint64 `uri:"u64"`
	}
	var s testStruct
	data := map[string]string{
		"u":   "1",
		"u8":  "255",
		"u16": "65535",
		"u32": "4294967295",
		"u64": "18446744073709551615",
	}
	require.NoError(t, Unmarshal(data, &s))
	assert.Equal(t, uint(1), s.U)
	assert.Equal(t, uint8(255), s.U8)
	assert.Equal(t, uint16(65535), s.U16)
	assert.Equal(t, uint32(4294967295), s.U32)
	assert.Equal(t, uint64(18446744073709551615), s.U64)
}

func TestUnmarshal_FloatTypes(t *testing.T) {
	type testStruct struct {
		F32 float32 `uri:"f32"`
		F64 float64 `uri:"f64"`
	}
	var s testStruct
	data := map[string]string{"f32": "3.14", "f64": "2.718281828"}
	require.NoError(t, Unmarshal(data, &s))
	assert.InDelta(t, float32(3.14), s.F32, 0.001)
	assert.InDelta(t, 2.718281828, s.F64, 0.000000001)
}

func TestUnmarshal_BoolValues(t *testing.T) {
	type testStruct struct {
		A bool `uri:"a"`
		B bool `uri:"b"`
		C bool `uri:"c"`
	}
	var s testStruct
	data := map[string]string{"a": "true", "b": "false", "c": "1"}
	require.NoError(t, Unmarshal(data, &s))
	assert.True(t, s.A)
	assert.False(t, s.B)
	assert.True(t, s.C)
}

func TestUnmarshal_InvalidInt(t *testing.T) {
	type testStruct struct {
		ID int `uri:"id"`
	}
	var s testStruct
	err := Unmarshal(map[string]string{"id": "notanumber"}, &s)
	assert.Error(t, err)
}

func TestUnmarshal_PartialData(t *testing.T) {
	type testStruct struct {
		A string `uri:"a"`
		B string `uri:"b"`
		C string `uri:"c"`
	}
	var s testStruct
	s.C = "default"
	data := map[string]string{"a": "hello", "b": "world"}
	require.NoError(t, Unmarshal(data, &s))
	assert.Equal(t, "hello", s.A)
	assert.Equal(t, "world", s.B)
	assert.Equal(t, "default", s.C) // preserved
}

func TestMarshal_Basic(t *testing.T) {
	type testStruct struct {
		ID   int    `uri:"id"`
		Name string `uri:"name"`
		Age  int    `json:"age"` // no uri tag, skipped
	}
	s := testStruct{ID: 42, Name: "bob", Age: 30}
	result, err := Marshal(s)
	require.NoError(t, err)
	assert.Equal(t, "42", result["id"])
	assert.Equal(t, "bob", result["name"])
	assert.NotContains(t, result, "age")
}

func TestMarshal_Pointer(t *testing.T) {
	type testStruct struct {
		ID int `uri:"id"`
	}
	s := &testStruct{ID: 99}
	result, err := Marshal(s)
	require.NoError(t, err)
	assert.Equal(t, "99", result["id"])
}

func TestMarshal_Bool(t *testing.T) {
	type testStruct struct {
		Active bool `uri:"active"`
	}
	s := testStruct{Active: true}
	result, err := Marshal(s)
	require.NoError(t, err)
	assert.Equal(t, "true", result["active"])
}

func TestMarshal_Float(t *testing.T) {
	type testStruct struct {
		Val float64 `uri:"val"`
	}
	s := testStruct{Val: 3.14}
	result, err := Marshal(s)
	require.NoError(t, err)
	assert.Contains(t, result["val"], "3.14")
}
