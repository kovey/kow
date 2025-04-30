package rule

import "fmt"

const (
	rule_gt = "gt"
)

type Gt struct {
	*Base
}

func NewGt() *Gt {
	return &Gt{Base: NewBase(rule_gt, nil)}
}

func (g *Gt) Valid(key string, val any, params ...any) (bool, error) {
	if len(params) != 1 {
		return false, fmt.Errorf("params[%+v] of field[%s] format error", params, key)
	}

	if !canCompare(val) || !canCompare(params[0]) {
		return false, fmt.Errorf("params[%v] and val[%v] of field[%s] can not compare", params[0], val, key)
	}

	res, err := compare(val, params[0])
	if err != nil {
		return false, fmt.Errorf("params[%v] and val[%v] of field[%s] compare error: %s", params[0], val, key, err)
	}

	rr := res == 1
	if !rr {
		return false, fmt.Errorf("value[%v] of field[%s] not greather than value[%v]", val, key, params[0])
	}

	return rr, nil
}
