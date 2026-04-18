package feistel

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

type Feistel struct {
	key    []byte
	rounds int
}

func NewFeistel(key []byte) *Feistel {
	return &Feistel{
		key:    key,
		rounds: 8,
	}
}

func split(x uint64) (uint32, uint32) {
	return uint32(x >> 32), uint32(x & 0xFFFFFFFF)
}

func join(left, right uint32) uint64 {
	return (uint64(left) << 32) | uint64(right)
}

func (c *Feistel) roundFunction(right uint32, roundKey []byte) uint32 {
	h := hmac.New(sha256.New, c.key)
	h.Write(roundKey)
	binary.Write(h, binary.BigEndian, right)
	sum := h.Sum(nil)
	return binary.BigEndian.Uint32(sum[:4])
}

func (c *Feistel) deriveRoundKey(round int) []byte {
	return []byte(fmt.Sprintf("round-%d", round))
}

func (c *Feistel) Encrypt(x uint64) uint64 {
	left, right := split(x)
	for i := 0; i < c.rounds; i++ {
		roundKey := c.deriveRoundKey(i)
		f := c.roundFunction(right, roundKey)
		newLeft := right
		newRight := left ^ f
		left, right = newLeft, newRight
	}
	return join(left, right)
}

func (c *Feistel) Decrypt(y uint64) uint64 {
	left, right := split(y)
	for i := c.rounds - 1; i >= 0; i-- {
		roundKey := c.deriveRoundKey(i)
		f := c.roundFunction(left, roundKey)
		newRight := left
		newLeft := right ^ f
		left, right = newLeft, newRight
	}
	return join(left, right)
}

func (c *Feistel) FPEEncrypt(value, min, max int64) (int64, error) {
	if value < min || value > max {
		return 0, errors.New("value out of range")
	}
	domainSize := max - min + 1
	if domainSize == 1 {
		return value, nil
	}
	offset := value - min
	cipher := c.Encrypt(uint64(offset))
	offset = int64(cipher % uint64(domainSize))
	return offset + min, nil
}

func (c *Feistel) FPEDecrypt(value, min, max int64) (int64, error) {
	if value < min || value > max {
		return 0, errors.New("value out of range")
	}
	domainSize := max - min + 1
	if domainSize == 1 {
		return value, nil
	}
	offset := value - min
	plain := c.Decrypt(uint64(offset))
	offset = int64(plain % uint64(domainSize))
	return offset + min, nil
}
