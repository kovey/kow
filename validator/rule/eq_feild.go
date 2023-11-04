package rule

const (
	rule_eq_feild = "eq_feild"
)

type EqFeild struct {
	*Base
}

func NewEqFeild() *EqFeild {
	return &EqFeild{Base: NewBase(rule_eq_feild, nil)}
}

func (g *EqFeild) Valid(key string, val any, params ...any) bool {
	if len(params) != 1 {
		return false
	}

	return val == params[0]
}
