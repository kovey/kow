package jwt

import (
	"encoding/json"
	"time"
)

type Paylaod[T any] struct {
	Iss string `json:"iss"`
	Iat int64  `json:"Iat"`
	Exp int64  `json:"Exp"`
	Jti string `json:"jti"`
	Ext T      `json:"ext"`
}

func NewPlayload[T any]() *Paylaod[T] {
	return &Paylaod[T]{}
}

func (p *Paylaod[T]) Encode() string {
	buf, err := json.Marshal(p)
	if err != nil {
		return "{}"
	}

	return string(buf)
}

func (p *Paylaod[T]) Decode(data string) error {
	return json.Unmarshal([]byte(data), p)
}

func (p *Paylaod[T]) IsExpired() bool {
	return p.Exp < time.Now().Unix()
}

func (p *Paylaod[T]) IsValid() bool {
	return p.Iat <= time.Now().Unix()
}
