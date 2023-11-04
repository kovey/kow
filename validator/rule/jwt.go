package rule

import (
	"encoding/base64"
	"strings"

	"github.com/kovey/debug-go/debug"
)

const (
	rule_jwt = "jwt"
)

type Jwt struct {
	*Base
}

func NewJwt() *Jwt {
	return &Jwt{NewBase(rule_jwt, nil)}
}

func (j *Jwt) Valid(key string, val any, params ...any) bool {
	tmp, ok := val.(string)
	if !ok {
		return false
	}

	info := strings.Split(tmp, ".")
	if len(info) != 3 {
		return false
	}

	for _, str := range info {
		if _, err := base64.RawURLEncoding.DecodeString(str); err != nil {
			debug.Erro("jwt valid failure, error: %s", err)
			return false
		}
	}

	return true
}
