package feistel

import (
	"testing"
)

func TestFeistel(t *testing.T) {
	key := []byte("32-bit-secret-key-for-demo") // 任意长度密钥
	cipher := NewFeistel(key)
	encs := make(map[uint64]bool)
	for i := uint64(0); i <= 10000; i++ {
		enc := cipher.Encrypt(i)
		dec := cipher.Decrypt(enc)
		t.Logf("原始值: %d, 加密值: %d, 解密值: %d\n", i, enc, dec)
		if i != dec {
			break
		}
		if _, ok := encs[enc]; ok {
			t.Logf("加密值: %d 重复", enc)
			break
		}
		encs[enc] = true
	}
}
