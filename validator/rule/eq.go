package rule

const (
	rule_eq = "eq"
)

type Eq struct {
	*Base
}

func NewEq() *Eq {
	return &Eq{Base: NewBase(rule_eq, nil)}
}

func (g *Eq) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		return false
	}

	return val == params[0]
}
