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

func (g *Eq) Valid(key string, val any, params ...any) (bool, error) {
	if len(params) != 1 {
		return false, fmt.Errorf("params: %+v format error", params)
	}

	res := val == params[0]
	if !res {
		return res, fmt.Errorf("value[%v] of field[%s] not equal give[%v]", val, key, params[0])
	}
	return res, nil
}
