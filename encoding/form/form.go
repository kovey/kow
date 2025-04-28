package form

import "reflect"

var Marshal func(v any) ([]byte, error) = formMarshal
var Unmarshal func(data map[string][]string, v any) error = formUnmarshal

func Init(m func(v any) ([]byte, error), u func(data map[string][]string, v any) error) {
	Marshal = m
	Unmarshal = u
}

func formMarshal(v any) ([]byte, error) {
	ref := &Reflect{Type: reflect.TypeOf(v)}
	if ref.Type.Kind() == reflect.Ptr {
		ref.Type = ref.Type.Elem()
	}

	return ref.Encode(v)
}

func formUnmarshal(data map[string][]string, v any) error {
	ref := NewReflect(v)
	return ref.Parse(data, v)
}
