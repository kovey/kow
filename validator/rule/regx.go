package rule

import (
	"fmt"
	"regexp"

	"github.com/kovey/debug-go/debug"
)

type Regx struct {
	*Base
}

const (
	rule_regx = "regx"
)

func NewRegx() *Regx {
	return &Regx{Base: NewBase(rule_regx, nil)}
}

func (r *Regx) Valid(key string, val any, params ...any) (bool, error) {
	if len(params) != 1 {
		return false, fmt.Errorf("params: %+v format error", params)
	}

	tmp, ok := val.(string)
	if !ok {
		return false, fmt.Errorf("val: %v not string", val)
	}
	pattern, ok := params[0].(string)
	if !ok {
		return false, fmt.Errorf("pattern: %v not string", params[0])
	}

	ok, err := regexp.Match(pattern, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure, error: %s", err)
		r.err = err
		return false, err
	}

	if !ok {
		return false, fmt.Errorf("value[%s] of field[%s] regx pattern[%v] failure", tmp, key, params[0])
	}

	return ok, nil
}
