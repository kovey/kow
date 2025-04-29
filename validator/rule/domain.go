package rule

import (
	"fmt"
	"regexp"

	"github.com/kovey/debug-go/debug"
)

const (
	rule_domain = "domain"
	domain_reg  = `[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(/.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+/.?`
)

type Domain struct {
	*Base
}

func NewDomain() *Domain {
	return &Domain{NewBase(rule_domain, nil)}
}

func (e *Domain) Valid(key string, val any, params ...any) bool {
	tmp, ok := val.(string)
	if !ok {
		return false
	}

	ok, err := regexp.Match(domain_reg, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure in domain, error: %s", err)
		e.err = fmt.Errorf("value[%s] of field[%s] is not domain", tmp, key)
		return false
	}

	if !ok {
		e.err = fmt.Errorf("value[%s] of field[%s] is not domain", tmp, key)
	}

	return ok
}
