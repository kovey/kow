package rule

import "fmt"

type ParamInterface interface {
	ValidParams() map[string]any
	Clone() ParamInterface
}

type Rule struct {
	Func   string
	Params []any
}

func NewRule(f string, params []any) *Rule {
	return &Rule{Func: f, Params: params}
}

type RuleInterface interface {
	Name() string
	Valid(key string, val any, params ...any) bool
	Err() error
}

type CompareInterface interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type BaseType interface {
	CompareInterface | ~string
}

type Base struct {
	name string
	err  error
}

func NewBase(name string, err error) *Base {
	return &Base{name: name, err: err}
}

func (b *Base) Name() string {
	return b.name
}

func (b *Base) Err() error {
	return b.err
}

func canCompare(data any) bool {
	switch data.(type) {
	case uint, uint8, uint16, uint32, uint64, float32, float64, int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

func convert[T BaseType](data any) (T, error) {
	var t T
	switch tmp := data.(type) {
	case T:
		return tmp, nil
	default:
		return t, fmt.Errorf("convert failure")
	}
}

func compare(left, right any) (int, error) {
	switch tmp := left.(type) {
	case uint:
		r, err := convert[uint](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case uint8:
		r, err := convert[uint8](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case uint16:
		r, err := convert[uint16](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case uint32:
		r, err := convert[uint32](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case uint64:
		r, err := convert[uint64](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case float32:
		r, err := convert[float32](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case float64:
		r, err := convert[float64](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case int:
		r, err := convert[int](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case int8:
		r, err := convert[int8](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case int16:
		r, err := convert[int16](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case int32:
		r, err := convert[int32](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	case int64:
		r, err := convert[int64](right)
		if err != nil {
			return 0, err
		}

		return compareBy(tmp, r), nil
	default:
		return 0, fmt.Errorf("unkown type")
	}
}

func compareBy[T CompareInterface](left, right T) int {
	if left > right {
		return 1
	}

	if left == right {
		return 0
	}

	return -1
}
