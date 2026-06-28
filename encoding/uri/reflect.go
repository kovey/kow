package uri

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

var (
	ErrNotPointer = errors.New("uri: object not pointer")
	ErrNotStruct  = errors.New("uri: object not struct")
)

const uriTag = "uri"

// Reflect holds the reflected type information for URI unmarshaling.
type Reflect struct {
	Type reflect.Type
}

func newReflect(val any) *Reflect {
	if val == nil {
		return &Reflect{}
	}
	rt := reflect.TypeOf(val)
	if rt.Kind() != reflect.Pointer {
		return &Reflect{}
	}
	return &Reflect{Type: rt.Elem()}
}

func (r *Reflect) parse(data map[string]string, val any) error {
	if val == nil {
		return ErrNotPointer
	}
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Pointer {
		return ErrNotPointer
	}
	if r.Type == nil {
		return ErrNotPointer
	}
	if r.Type.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	v := rv.Elem()
	for i := 0; i < r.Type.NumField(); i++ {
		field := r.Type.Field(i)
		tag := field.Tag.Get(uriTag)
		if tag == "" || tag == "-" {
			continue
		}

		paramVal, ok := data[tag]
		if !ok {
			continue
		}

		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		if err := setField(fv, field.Type, paramVal); err != nil {
			return err
		}
	}
	return nil
}

func setField(fv reflect.Value, ft reflect.Type, val string) error {
	// Check for custom Unmarshaler first
	if fv.CanAddr() {
		if u, ok := fv.Addr().Interface().(Unmarshaler); ok {
			return u.UnmarshalURI(val)
		}
	}

	// Handle time.Time
	if _, ok := fv.Interface().(time.Time); ok {
		t, err := time.Parse(time.DateTime, val)
		if err != nil {
			return err
		}
		fv.Set(reflect.ValueOf(t))
		return nil
	}

	switch ft.Kind() {
	case reflect.String:
		fv.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 10, ft.Bits())
		if err != nil {
			return err
		}
		fv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(val, 10, ft.Bits())
		if err != nil {
			return err
		}
		fv.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(val, ft.Bits())
		if err != nil {
			return err
		}
		fv.SetFloat(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	default:
		// unsupported type, skip silently
	}
	return nil
}

func (r *Reflect) encode(val any) (map[string]string, error) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	result := make(map[string]string)
	for i := 0; i < r.Type.NumField(); i++ {
		field := r.Type.Field(i)
		tag := field.Tag.Get(uriTag)
		if tag == "" || tag == "-" {
			continue
		}

		fv := v.Field(i)
		s, err := fieldToString(fv, field.Type)
		if err != nil {
			return nil, err
		}
		result[tag] = s
	}
	return result, nil
}

func fieldToString(fv reflect.Value, ft reflect.Type) (string, error) {
	switch ft.Kind() {
	case reflect.String:
		return fv.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(fv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(fv.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(fv.Float(), 'g', -1, ft.Bits()), nil
	case reflect.Bool:
		return strconv.FormatBool(fv.Bool()), nil
	default:
		return "", nil
	}
}
