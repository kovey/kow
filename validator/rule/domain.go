package rule

import (
	"fmt"
	"regexp"

	"github.com/kovey/debug-go/debug"
)

const (
	rule_domain = "domain"
	domain_reg  = `^([a-zA-Z0-9][a-zA-Z0-9\-]{1,61}[a-zA-Z0-9]\.)+[a-zA-Z0-9]{2,6}$`
)

type Domain struct {
	*Base
}

func NewDomain() *Domain {
	return &Domain{NewBase(rule_domain, nil)}
}

func (e *Domain) Valid(key string, val any, params ...any) (bool, error) {
	tmp, ok := val.(string)
	if !ok {
		return false, e.err
	}

	ok, err := regexp.Match(domain_reg, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure in domain, error: %s", err)
		return false, fmt.Errorf("value[%s] of field[%s] is not domain", tmp, key)
	}

	if !ok {
		return ok, fmt.Errorf("value[%s] of field[%s] is not domain", tmp, key)
	}

	return ok, nil
}
