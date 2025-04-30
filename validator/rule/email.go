package rule

import (
	"fmt"
	"regexp"

	"github.com/kovey/debug-go/debug"
)

const (
	rule_email = "email"
	email_reg  = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
)

type Email struct {
	*Base
}

func NewEmail() *Email {
	return &Email{NewBase(rule_email, nil)}
}

func (e *Email) Valid(key string, val any, params ...any) (bool, error) {
	tmp, ok := val.(string)
	if !ok {
		return false, e.err
	}

	ok, err := regexp.Match(email_reg, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure in email, error: %s", err)
		return false, fmt.Errorf("value[%s] of field[%s] is not email", tmp, key)
	}

	if !ok {
		return ok, fmt.Errorf("value[%s] of field[%s] is not email", tmp, key)
	}
	return ok, nil
}
