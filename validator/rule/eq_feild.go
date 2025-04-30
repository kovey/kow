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

func (g *EqFeild) Valid(key string, val any, params ...any) (bool, error) {
	if len(params) != 1 {
		return false, fmt.Errorf("params: %+v format error", params)
	}

	res := val == params[0]
	if !res {
		return res, fmt.Errorf("value[%v] of field[%s] not equal value[%v]", val, key, params[0])
	}
	return res, nil
}
