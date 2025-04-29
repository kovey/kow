package form

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Err_Not_Pointer = errors.New("object not pointer")
var Err_Not_Struct = errors.New("object not struct")

const (
	form_tag = "form"
)

type Reflect struct {
	Type reflect.Type
}

func NewReflect(val any) *Reflect {
	return &Reflect{Type: reflect.TypeOf(val).Elem()}
}

func (f *Reflect) _parse(data map[string][]string, val any) error {
	v := reflect.ValueOf(val).Elem()
	for i := 0; i < f.Type.NumField(); i++ {
		field := f.Type.Field(i)
		dt, ok := data[field.Tag.Get(form_tag)]
		if !ok || len(dt) == 0 {
			continue
		}

		vField := v.Field(i)
		orData := dt[0]
		if _, ok := vField.Interface().(time.Time); ok {
			tmpTime, err := time.Parse(time.DateTime, orData)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(tmpTime))
		}
		switch field.Type.Kind() {
		case reflect.String:
			vField.Set(reflect.ValueOf(orData))
		case reflect.Int:
			tmpVal, err := strconv.Atoi(orData)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(tmpVal))
		case reflect.Int8:
			tmpVal, err := strconv.Atoi(orData)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(int8(tmpVal)))
		case reflect.Int16:
			tmpVal, err := strconv.Atoi(orData)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(int16(tmpVal)))
		case reflect.Int32:
			tmpVal, err := strconv.Atoi(orData)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(int32(tmpVal)))
		case reflect.Int64:
			tmpVal, err := strconv.ParseInt(orData, 10, 64)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(tmpVal))
		case reflect.Uint:
			tmpVal, err := strconv.ParseUint(orData, 10, 64)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(uint(tmpVal)))
		case reflect.Uint8:
			tmpVal, err := strconv.Atoi(orData)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(uint8(tmpVal)))
		case reflect.Uint16:
			tmpVal, err := strconv.Atoi(orData)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(uint16(tmpVal)))
		case reflect.Uint32:
			tmpVal, err := strconv.ParseUint(orData, 10, 64)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(uint32(tmpVal)))
		case reflect.Uint64:
			tmpVal, err := strconv.ParseUint(orData, 10, 64)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(tmpVal))
		case reflect.Bool:
			tmpVal, err := strconv.ParseBool(orData)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(tmpVal))
		case reflect.Float32:
			tmpVal, err := strconv.ParseFloat(orData, 10)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(float32(tmpVal)))
		case reflect.Float64:
			tmpVal, err := strconv.ParseFloat(orData, 10)
			if err != nil {
				return err
			}

			vField.Set(reflect.ValueOf(tmpVal))
		case reflect.Array, reflect.Slice:
			vTypeValue := reflect.New(field.Type.Elem()).Elem()
			if _, ok := vTypeValue.Interface().(time.Time); ok {
				for _, orData := range dt {
					tmpTime, err := time.Parse(time.DateTime, orData)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(tmpTime)))
				}
			}
			switch field.Type.Elem().Kind() {
			case reflect.String:
				for _, orData := range dt {
					vField.Set(reflect.Append(vField, reflect.ValueOf(orData)))
				}
			case reflect.Int:
				for _, orData := range dt {
					tmpVal, err := strconv.Atoi(orData)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(tmpVal)))
				}
			case reflect.Int8:
				for _, orData := range dt {
					tmpVal, err := strconv.Atoi(orData)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(int8(tmpVal))))
				}
			case reflect.Int16:
				for _, orData := range dt {
					tmpVal, err := strconv.Atoi(orData)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(int16(tmpVal))))
				}
			case reflect.Int32:
				for _, orData := range dt {
					tmpVal, err := strconv.Atoi(orData)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(int32(tmpVal))))
				}
			case reflect.Int64:
				for _, orData := range dt {
					tmpVal, err := strconv.ParseInt(orData, 10, 64)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(tmpVal)))
				}
			case reflect.Uint:
				for _, orData := range dt {
					tmpVal, err := strconv.ParseUint(orData, 10, 64)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(uint(tmpVal))))
				}
			case reflect.Uint8:
				for _, orData := range dt {
					tmpVal, err := strconv.Atoi(orData)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(uint8(tmpVal))))
				}
			case reflect.Uint16:
				for _, orData := range dt {
					tmpVal, err := strconv.Atoi(orData)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(uint16(tmpVal))))
				}
			case reflect.Uint32:
				for _, orData := range dt {
					tmpVal, err := strconv.ParseUint(orData, 10, 64)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(uint32(tmpVal))))
				}
			case reflect.Uint64:
				for _, orData := range dt {
					tmpVal, err := strconv.ParseUint(orData, 10, 64)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(tmpVal)))
				}
			case reflect.Bool:
				for _, orData := range dt {
					tmpVal, err := strconv.ParseBool(orData)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(tmpVal)))
				}
			case reflect.Float32:
				for _, orData := range dt {
					tmpVal, err := strconv.ParseFloat(orData, 10)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(float32(tmpVal))))
				}
			case reflect.Float64:
				for _, orData := range dt {
					tmpVal, err := strconv.ParseFloat(orData, 10)
					if err != nil {
						return err
					}

					vField.Set(reflect.Append(vField, reflect.ValueOf(tmpVal)))
				}
			}
		}
	}

	return nil
}

func (f *Reflect) Parse(data map[string][]string, val any) error {
	if f.Type == nil {
		return Err_Not_Pointer
	}

	if f.Type.Kind() != reflect.Struct {
		return Err_Not_Struct
	}

	return f._parse(data, val)
}

func (f *Reflect) Encode(val any) ([]byte, error) {
	var builder strings.Builder
	vVal := reflect.ValueOf(val)
	if vVal.Kind() == reflect.Pointer {
		vVal = vVal.Elem()
	}
	for i := 0; i < f.Type.NumField(); i++ {
		if i > 0 {
			builder.WriteString("&")
		}
		field := f.Type.Field(i)
		name := field.Tag.Get(form_tag)
		if name == "" {
			name = field.Name
		}
		vField := vVal.Field(i)
		switch field.Type.Kind() {
		case reflect.Struct:
			if _, ok := vField.Interface().(time.Time); ok {
				data := vField.Interface().(time.Time)
				builder.WriteString(name)
				builder.WriteString("=")
				builder.WriteString(data.Format(time.DateTime))
			}
		case reflect.String:
			builder.WriteString(name)
			builder.WriteString("=")
			builder.WriteString(vField.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			builder.WriteString(name)
			builder.WriteString("=")
			builder.WriteString(strconv.FormatInt(vField.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			builder.WriteString(name)
			builder.WriteString("=")
			builder.WriteString(strconv.FormatUint(vField.Uint(), 10))
		case reflect.Bool:
			builder.WriteString(name)
			builder.WriteString("=")
			builder.WriteString(strconv.FormatBool(vField.Bool()))
		case reflect.Float32, reflect.Float64:
			builder.WriteString(name)
			builder.WriteString("=")
			builder.WriteString(strconv.FormatFloat(vField.Float(), 'g', 10, 64))
		case reflect.Array, reflect.Slice:
			switch field.Type.Elem().Kind() {
			case reflect.Struct:
				for i := 0; i < vField.Len(); i++ {
					data, ok := vField.Index(i).Interface().(time.Time)
					if ok {
						builder.WriteString(name)
						builder.WriteString("[]=")
						builder.WriteString(data.Format(time.DateTime))
					}
				}
			case reflect.String:
				for i := 0; i < vField.Len(); i++ {
					builder.WriteString(name)
					builder.WriteString("[]=")
					builder.WriteString(vField.Index(i).String())
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				for i := 0; i < vField.Len(); i++ {
					builder.WriteString(name)
					builder.WriteString("[]=")
					builder.WriteString(strconv.FormatInt(vField.Index(i).Int(), 10))
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				for i := 0; i < vField.Len(); i++ {
					builder.WriteString(name)
					builder.WriteString("[]=")
					builder.WriteString(strconv.FormatUint(vField.Index(i).Uint(), 10))
				}
			case reflect.Bool:
				for i := 0; i < vField.Len(); i++ {
					builder.WriteString(name)
					builder.WriteString("[]=")
					builder.WriteString(strconv.FormatBool(vField.Index(i).Bool()))
				}
			case reflect.Float32, reflect.Float64:
				for i := 0; i < vField.Len(); i++ {
					builder.WriteString(name)
					builder.WriteString("[]=")
					builder.WriteString(strconv.FormatFloat(vField.Index(i).Float(), 'g', 10, 64))
				}
			}
		}
	}

	return []byte(builder.String()), nil
}
