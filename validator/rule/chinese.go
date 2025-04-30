package rule

import (
	"fmt"
	"regexp"

	"github.com/kovey/debug-go/debug"
)

const (
	rule_chinese = "chinese"
	chinese_reg  = "^[\u4e00-\u9fa5]+$"
)

type Chinese struct {
	*Base
}

func NewChinese() *Chinese {
	return &Chinese{NewBase(rule_chinese, nil)}
}

func (c *Chinese) Valid(key string, val any, params ...any) (bool, error) {
	tmp, ok := val.(string)
	if !ok {
		return false, c.err
	}

	ok, err := regexp.Match(chinese_reg, []byte(tmp))
	if err != nil {
		debug.Erro("regexp matched failure in chinese, error: %s", err)
		return false, fmt.Errorf("value[%s] of field[%s] is not chinese", tmp, key)
	}

	if !ok {
		return ok, fmt.Errorf("value[%s] of field[%s] is not chinese", tmp, key)
	}

	return ok, nil
}
