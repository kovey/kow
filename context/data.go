package context

import (
	"strconv"

	"github.com/kovey/debug-go/debug"
)

type Params map[string]string

func (p Params) Reset() {
	for k := range p {
		delete(p, k)
	}
}

func (p Params) GetString(key string) string {
	if val, ok := p[key]; ok {
		return val
	}

	return ""
}

func (p Params) GetInt(key string) int {
	if val, ok := p[key]; ok {
		tmp, err := strconv.Atoi(val)
		if err != nil {
			debug.Erro(err.Error())
			return 0
		}

		return tmp
	}

	return 0
}

func (p Params) GetFloat(key string) float64 {
	if val, ok := p[key]; ok {
		tmp, err := strconv.ParseFloat(val, 64)
		if err != nil {
			debug.Erro(err.Error())
			return 0
		}

		return tmp
	}

	return 0
}

func (p Params) GetBool(key string) bool {
	if val, ok := p[key]; ok {
		tmp, err := strconv.ParseBool(val)
		if err != nil {
			debug.Erro(err.Error())
			return false
		}

		return tmp
	}

	return false
}

type Data map[string]any
