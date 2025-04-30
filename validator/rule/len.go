package rule

import "fmt"

const (
	rule_len = "len"
)

type Len struct {
	*Base
}

func NewLen() *Len {
	return &Len{Base: NewBase(rule_len, nil)}
}

func (g *Len) Valid(key string, val any, params ...any) (bool, error) {
	if len(params) != 1 {
		return false, fmt.Errorf("params[%+v] of field[%s] format error", params, key)
	}

	v, err := convert[string](val)
	if err != nil {
		return false, err
	}

	l, err := convert[int](params[0])
	if err != nil {
		return false, err
	}

	rr := len(v) == l
	if !rr {
		return false, fmt.Errorf("value[%v] of field[%s] length not equal value[%v]", val, key, params[0])
	}

	return true, nil
}
