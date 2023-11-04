package rule

import (
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

func (u *Url) Valid(key string, val any, params ...any) bool {
	tmp, ok := val.(string)
	if !ok {
		return false
	}

	ok, err := regexp.Match(url_reg, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure in url, error: %s", err)
		return false
	}

	return ok
}
