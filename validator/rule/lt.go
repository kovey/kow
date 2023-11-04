package rule

const (
	rule_lt = "lt"
)

type Lt struct {
	*Base
}

func NewLt() *Lt {
	return &Lt{Base: NewBase(rule_lt, nil)}
}

func (g *Lt) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		return false
	}

	if !canCompare(val) || !canCompare(params[0]) {
		return false
	}

	res, err := compare(val, params[0])
	if err != nil {
		return false
	}

	return res == -1
}
