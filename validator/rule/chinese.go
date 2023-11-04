package rule

import (
	"regexp"

	"github.com/kovey/debug-go/debug"
)

const (
	rule_chinese = "chinese"
	chinese_reg  = `[\u4e00-\u9fa5]`
)

type Chinese struct {
	*Base
}

func NewChinese() *Chinese {
	return &Chinese{NewBase(rule_chinese, nil)}
}

func (c *Chinese) Valid(key string, val any, params ...any) bool {
	tmp, ok := val.(string)
	if !ok {
		return false
	}

	ok, err := regexp.Match(chinese_reg, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure in chinese, error: %s", err)
		return false
	}

	return ok
}
