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

func (r *Regx) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		r.err = fmt.Errorf("params: %+v format error", params)
		return false
	}

	tmp, ok := val.(string)
	if !ok {
		r.err = fmt.Errorf("val: %v not string", val)
		return false
	}
	pattern, ok := params[0].(string)
	if !ok {
		r.err = fmt.Errorf("pattern: %v not string", params[0])
		return false
	}

	ok, err := regexp.Match(pattern, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure, error: %s", err)
		r.err = err
		return false
	}

	return ok
}
