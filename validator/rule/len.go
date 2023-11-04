package rule

const (
	rule_len = "len"
)

type Len struct {
	*Base
}

func NewLen() *Len {
	return &Len{Base: NewBase(rule_len, nil)}
}

func (g *Len) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		return false
	}

	v, err := convert[string](val)
	if err != nil {
		return false
	}

	l, err := convert[int](params[0])
	if err != nil {
		return false
	}
	return len(v) == l
}
