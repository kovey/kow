package rule

import (
	"encoding/base64"
	"fmt"
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

func (j *Jwt) Valid(key string, val any, params ...any) (bool, error) {
	tmp, ok := val.(string)
	if !ok {
		if j.err != nil {
			return false, j.err
		}

		return false, fmt.Errorf("val is not string")
	}

	info := strings.Split(tmp, ".")
	if len(info) != 3 {
		return false, fmt.Errorf("value[%s] of field[%s] is not jwt", tmp, key)
	}

	for _, str := range info {
		if _, err := base64.RawURLEncoding.DecodeString(str); err != nil {
			debug.Erro("jwt valid failure, error: %s", err)
			return false, fmt.Errorf("value[%s] of field[%s] is not jwt", tmp, key)
		}
	}

	return true, nil
}
