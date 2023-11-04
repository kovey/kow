package rule

const (
	rule_max = "maxlen"
)

type MaxLen struct {
	*Base
}

func NewMaxLen() *MaxLen {
	return &MaxLen{Base: NewBase(rule_max, nil)}
}

func (m *MaxLen) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		return false
	}

	tmp := params[0].(int)
	switch vv := val.(type) {
	case string:
		return len(vv) <= tmp
	}

	return false
}
