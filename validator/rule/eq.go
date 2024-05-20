package rule

import "fmt"

const (
	rule_eq = "eq"
)

type Eq struct {
	*Base
}

func NewEq() *Eq {
	return &Eq{Base: NewBase(rule_eq, nil)}
}

func (g *Eq) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		g.err = fmt.Errorf("params: %+v format error", params)
		return false
	}

	res := val == params[0]
	if !res {
		g.err = fmt.Errorf("value[%v] of field[%s] not equal give[%v]", val, key, params[0])
	}
	return res
}
