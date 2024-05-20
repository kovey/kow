package validator

import (
	"fmt"

	"github.com/kovey/kow/validator/rule"
)

type Validator struct {
	data map[string]rule.RuleInterface
}

func NewValidator() *Validator {
	return &Validator{data: make(map[string]rule.RuleInterface)}
}

func (v *Validator) Add(valid rule.RuleInterface) {
	v.data[valid.Name()] = valid
}

func (v *Validator) Valid(key string, value any, rules []*rule.Rule, all map[string]any) error {
	for _, rule := range rules {
		valid, ok := v.data[rule.Func]
		if !ok {
			return fmt.Errorf("validator[%s] not found", rule.Func)
		}

		if rule.Func == "eq_feild" {
			if len(rule.Params) != 1 {
				return fmt.Errorf("param[%s] valid with[%s] failure, params count not 1", key, rule.Func)
			}
			kk, ok := rule.Params[0].(string)
			if !ok {
				return fmt.Errorf("param[%s] valid with[%s] failure, params[%v] not string", key, rule.Func, rule.Params[0])
			}
			val, ok := all[kk]
			if !ok {
				return fmt.Errorf("param[%s] valid with[%s] failure, field[%s] not found", key, rule.Func, kk)
			}
			if valid.Valid(key, value, val) {
				continue
			}
		} else {
			if valid.Valid(key, value, rule.Params...) {
				continue
			}
		}

		if err := valid.Err(); err != nil {
			return err
		}

		return fmt.Errorf("param[%s] valid with[%s] failure", key, rule.Func)
	}

	return nil
}

var v = NewValidator()
var r = newParamRules()

func init() {
	v.Add(rule.NewEq())
	v.Add(rule.NewEqFeild())
	v.Add(rule.NewGe())
	v.Add(rule.NewGt())
	v.Add(rule.NewLe())
	v.Add(rule.NewLen())
	v.Add(rule.NewLt())
	v.Add(rule.NewMaxLen())
	v.Add(rule.NewMinLen())
	v.Add(rule.NewRegx())
	v.Add(rule.NewUrl())
	v.Add(rule.NewChinese())
	v.Add(rule.NewDomain())
	v.Add(rule.NewEmail())
	v.Add(rule.NewJwt())
}

// rule format:
// eq:int:1,2 (eq is rule name, int is param type, 1, 2 are params)
//
// rule names:
// eq,eq_feild,ge,gt,le,len,lt,minlen,maxlen,regx,chinese,email,jwt,url
func RegRule(key string, rules ...string) bool {
	return r.add(key, rules...)
}

func Register(r rule.RuleInterface) {
	v.Add(r)
}

func Valid(all map[string]any) error {
	for key, value := range all {
		if err := v.Valid(key, value, r.get(key), all); err != nil {
			return err
		}
	}

	return nil
}
