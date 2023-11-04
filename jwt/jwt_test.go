package jwt

import (
	"fmt"
	"testing"
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
	if err != nil {
		t.Fatalf("decode failure: %s", err)
	}

	fmt.Println("token: ", tk)
	ex, err := jwt.Decode(tk)
	fmt.Printf("ex: %+v, err: %s\n", ex, err)
}
