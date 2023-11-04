package rule

const (
	rule_gt = "gt"
)

type Gt struct {
	*Base
}

func NewGt() *Gt {
	return &Gt{Base: NewBase(rule_gt, nil)}
}

func (g *Gt) Valid(key string, val any, params ...any) bool {
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

	return res == 1
}
