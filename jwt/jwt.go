package jwt

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kovey/kow/encoding/json"
)

const (
	Jwt_Alg        = "HS256"
	Jwt_Type       = "JWT"
	Jwt_Admin      = "JWT_API_ADMIN_KOVEY"
	Jwt_Alg_Sha256 = "sha256"
)

var Err_Token_Format_Invalid = errors.New("token format error")
var Err_Token_Sign_Invalid = errors.New("token sign error")
var Err_Token_Expired = errors.New("token expired")
var Err_Token_Iat_Invalid = errors.New("token iat invalid")

type Jwt[T any] struct {
	Expired   int32
	Header    Header
	Key       string
	AlgConfig Header
}

func NewJwt[T any](key string, expired int32) *Jwt[T] {
	j := &Jwt[T]{Key: key, Expired: expired, Header: make(Header), AlgConfig: make(Header)}
	j.AlgConfig.Add(Jwt_Alg, Jwt_Alg_Sha256)
	j.Header.Add("alg", Jwt_Alg)
	j.Header.Add("typ", Jwt_Type)
	return j
}

func (j *Jwt[T]) Encode(ext T) (string, error) {
	hStr, err := json.Marshal(j.Header)
	if err != nil {
		return "", err
	}

	header, err := Encrypt(j.Key, string(hStr))
	if err != nil {
		return "", err
	}

	load := NewPlayload[T]()
	load.Iss = Jwt_Admin
	now := time.Now()
	load.Iat = now.Unix()
	load.Exp = now.Unix() + int64(j.Expired)
	load.Jti = Sha256(fmt.Sprintf("%s-%d-%d", Jwt_Admin, now.UnixNano(), rand.Int63()))
	load.Ext = ext
	payload, err := Encrypt(j.Key, load.Encode())
	if err != nil {
		return "", err
	}

	sign, err := j.signature(fmt.Sprintf("%s.%s", header, payload), j.Key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s.%s", header, payload, sign), nil
}

func (j *Jwt[T]) Decode(token string) (T, error) {
	var res T
	info := strings.Split(token, ".")
	if len(info) != 3 {
		return res, Err_Token_Format_Invalid
	}

	tmp, err := Decrypt(j.Key, info[0])
	if err != nil {
		return res, err
	}

	header := make(Header)
	if err := json.Unmarshal([]byte(tmp), &header); err != nil {
		return res, err
	}

	if header.Get("alg") != Jwt_Alg || header.Get("typ") != Jwt_Type {
		return res, Err_Token_Format_Invalid
	}

	sign, err := j.signature(fmt.Sprintf("%s.%s", info[0], info[1]), j.Key)
	if err != nil {
		return res, err
	}

	if sign != info[2] {
		return res, Err_Token_Sign_Invalid
	}

	payload := NewPlayload[T]()
	p, err := Decrypt(j.Key, info[1])
	if err != nil {
		return res, err
	}

	if err := payload.Decode(p); err != nil {
		return res, err
	}

	if payload.IsExpired() {
		return res, Err_Token_Expired
	}

	if !payload.IsValid() {
		return res, Err_Token_Iat_Invalid
	}

	return payload.Ext, nil
}

func (j *Jwt[T]) signature(input, key string) (string, error) {
	return Encrypt(key, HmacSha256(input, key))
}
