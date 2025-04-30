package rule

import (
	"fmt"
	"regexp"

	"github.com/kovey/debug-go/debug"
)

const (
	rule_url = "url"
	url_reg  = `^(http|https)://[a-zA-Z0-9]+(.[a-zA-Z0-9]+)+([a-zA-Z0-9-._?,'+/\~:#[]@!$&*])*$`
)

type Url struct {
	*Base
}

func NewUrl() *Url {
	return &Url{NewBase(rule_url, nil)}
}

func (u *Url) Valid(key string, val any, params ...any) (bool, error) {
	tmp, ok := val.(string)
	if !ok {
		return false, fmt.Errorf("value[%v] of field[%s] not string", val, key)
	}

	ok, err := regexp.Match(url_reg, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure in url, error: %s", err)
		return false, fmt.Errorf("value[%s] of field[%s] is not url", tmp, key)
	}

	if !ok {
		return ok, fmt.Errorf("value[%s] of field[%s] is not url", tmp, key)
	}

	return ok, nil
}
