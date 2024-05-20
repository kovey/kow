package rule

import "fmt"

const (
	rule_ge = "ge"
)

type Ge struct {
	*Base
}

func NewGe() *Ge {
	return &Ge{Base: NewBase(rule_ge, nil)}
}

func (g *Ge) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		g.err = fmt.Errorf("params[%+v] of field[%s] format error", params, key)
		return false
	}

	if !canCompare(val) || !canCompare(params[0]) {
		g.err = fmt.Errorf("params[%v] and val[%v] of field[%s] can not compare", params[0], val, key)
		return false
	}

	res, err := compare(val, params[0])
	if err != nil {
		g.err = fmt.Errorf("params[%v] and val[%v] of field[%s] compare error: %s", params[0], val, key, err)
		return false
	}

	rr := res >= 0
	if !rr {
		g.err = fmt.Errorf("value[%v] of field[%s] not greather than or equal value[%v]", val, key, params[0])
	}

	return rr
}
