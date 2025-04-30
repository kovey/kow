package jwt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Ext struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
}

func TestEncode(t *testing.T) {
	key := "vGKUOiH8jF6z9atNR3Ty3po4rVXQV1Qa9UzNV91mO9f"
	jwt := NewJwt[Ext](key, 80000)

	ext := Ext{}
	ext.UserId = 1
	ext.Name = "aaa"
	tk, err := jwt.Encode(ext)
	assert.Nil(t, err)
	assert.True(t, len(strings.Split(tk, ".")) == 3)
	ex, err := jwt.Decode(tk)
	assert.Nil(t, err)
	assert.Equal(t, ext, ex)
}
