package rule

const (
	rule_le = "le"
)

type Le struct {
	*Base
}

func NewLe() *Le {
	return &Le{Base: NewBase(rule_le, nil)}
}

func (g *Le) Valid(key string, val any, params ...any) bool {
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

	return res <= 0
}
