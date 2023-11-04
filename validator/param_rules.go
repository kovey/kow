package validator

import (
	"strconv"
	"strings"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/validator/rule"
)

type paramRules struct {
	rules map[string][]*rule.Rule
}

func newParamRules() *paramRules {
	return &paramRules{rules: make(map[string][]*rule.Rule)}
}

// rule: eq:int:1,2
func (p *paramRules) add(key string, rules ...string) bool {
	if _, ok := p.rules[key]; ok {
		return false
	}

	p.rules[key] = make([]*rule.Rule, len(rules))
	for i, r := range rules {
		rr := parseRule(r)
		if rr == nil {
			return false
		}

		p.rules[key][i] = rr
	}

	return true
}

func (p *paramRules) get(key string) []*rule.Rule {
	if rules, ok := p.rules[key]; ok {
		return rules
	}

	return nil
}

func parseRule(ru string) *rule.Rule {
	info := strings.Split(ru, ":")
	switch len(info) {
	case 1:
		return rule.NewRule(ru, nil)
	case 2:
		debug.Erro("rule[%s] format error", ru)
		return nil
	default:
		params := strings.Split(info[2], ",")
		tmp := make([]any, len(params))
		switch info[1] {
		case "string":
			for i, p := range params {
				tmp[i] = p
			}
		case "int":
			for i, p := range params {
				pa, err := strconv.Atoi(p)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = pa
			}
		case "int8":
			for i, p := range params {
				pa, err := strconv.ParseInt(p, 10, 8)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = int8(pa)
			}
		case "int16":
			for i, p := range params {
				pa, err := strconv.ParseInt(p, 10, 16)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = int16(pa)
			}
		case "int32":
			for i, p := range params {
				pa, err := strconv.ParseInt(p, 10, 32)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = int32(pa)
			}
		case "int64":
			for i, p := range params {
				pa, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = pa
			}
		case "float32":
			for i, p := range params {
				pa, err := strconv.ParseFloat(p, 32)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = float32(pa)
			}
		case "float64":
			for i, p := range params {
				pa, err := strconv.ParseFloat(p, 64)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = pa
			}
		case "uint":
			for i, p := range params {
				pa, err := strconv.ParseUint(p, 10, 64)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = pa
			}
		case "uint8":
			for i, p := range params {
				pa, err := strconv.ParseUint(p, 10, 8)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = int8(pa)
			}
		case "uint16":
			for i, p := range params {
				pa, err := strconv.ParseUint(p, 10, 16)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = int16(pa)
			}
		case "uint32":
			for i, p := range params {
				pa, err := strconv.ParseUint(p, 10, 32)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = int32(pa)
			}
		case "uint64":
			for i, p := range params {
				pa, err := strconv.ParseUint(p, 10, 64)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = pa
			}
		case "bool":
			for i, p := range params {
				pa, err := strconv.ParseBool(p)
				if err != nil {
					debug.Erro("rule[%s] params[%d] value[%s] is not[%s]", info[0], i, p, info[2])
					return nil
				}
				tmp[i] = pa
			}
		default:
			return nil
		}

		return rule.NewRule(info[0], tmp)
	}
}
