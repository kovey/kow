package rule

import "fmt"

const (
	rule_eq_feild = "eq_feild"
)

type EqFeild struct {
	*Base
}

func NewEqFeild() *EqFeild {
	return &EqFeild{Base: NewBase(rule_eq_feild, nil)}
}

func (g *EqFeild) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		g.err = fmt.Errorf("params: %+v format error", params)
		return false
	}

	res := val == params[0]
	if !res {
		g.err = fmt.Errorf("value[%v] of field[%s] not equal value[%v]", val, key, params[0])
	}
	return res
}
