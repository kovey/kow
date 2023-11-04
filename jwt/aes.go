package jwt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/kovey/debug-go/debug"
)

func Encrypt(key string, data string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(key + "=")
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(res)
	if err != nil {
		return "", err
	}

	size := block.BlockSize()
	tmp := pkcs7Padding([]byte(data), size)
	iv := res[:size]
	mode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(tmp))
	mode.CryptBlocks(crypted, tmp)

	return base64.RawURLEncoding.EncodeToString(crypted), nil
}

func Decrypt(key, data string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(key + "=")
	if err != nil {
		debug.Erro("base64 decode key[%s] failure, error: %s", key, err)
		return "", err
	}
	tmp, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		debug.Erro("base64 decode data[%s] failure, error: %s", data, err)
		return "", err
	}

	block, err := aes.NewCipher(res)
	if err != nil {
		return "", err
	}

	size := block.BlockSize()
	iv := res[:size]
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(tmp))
	mode.CryptBlocks(decrypted, tmp)

	decrypted = pkcs7UnPadding(decrypted)
	return string(decrypted), nil
}

func Sha256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return string(h.Sum(nil))
}

func HmacSha256(data string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return string(mac.Sum(nil))
}
func pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length <= unpadding {
		return nil
	}

	return origData[:(length - unpadding)]
}
