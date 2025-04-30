package rule

import "fmt"

const (
	rule_max = "maxlen"
)

type MaxLen struct {
	*Base
}

func NewMaxLen() *MaxLen {
	return &MaxLen{Base: NewBase(rule_max, nil)}
}

func (m *MaxLen) Valid(key string, val any, params ...any) (bool, error) {
	if len(params) != 1 {
		return false, fmt.Errorf("params[%+v] of field[%s] format error", params, key)
	}

	tmp := params[0].(int)
	switch vv := val.(type) {
	case string:
		rr := len(vv) <= tmp
		if !rr {
			return rr, fmt.Errorf("value[%v] of field[%s] length out range of max value[%v]", val, key, params[0])
		}

		return rr, nil
	}

	return false, fmt.Errorf("value[%v] of field[%s] not string", val, key)
}
