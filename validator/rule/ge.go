package rule

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
		return false
	}

	if !canCompare(val) || !canCompare(params[0]) {
		return false
	}

	res, err := compare(val, params[0])
	if err != nil {
		return false
	}

	return res >= 0
}
