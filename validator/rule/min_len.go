package rule

const (
	rule_min = "minlen"
)

type MinLen struct {
	*Base
}

func NewMinLen() *MinLen {
	return &MinLen{Base: NewBase(rule_min, nil)}
}

func (m *MinLen) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		return false
	}

	tmp := params[0].(int)
	switch vv := val.(type) {
	case string:
		return len(vv) >= tmp
	}

	return false
}
