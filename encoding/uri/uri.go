package uri

import "reflect"

// Unmarshaler is the interface implemented by types that can unmarshal
// a URI path parameter value into themselves.
type Unmarshaler interface {
	UnmarshalURI(data string) error
}

// Unmarshal parses URI path parameters (map[string]string) and stores the
// result in the struct pointed to by v. Fields are matched by the `uri` tag.
func Unmarshal(data map[string]string, v any) error {
	ref := newReflect(v)
	return ref.parse(data, v)
}

// Marshal encodes a struct into a map of URI parameters. Fields with a `uri`
// tag are included in the output.
func Marshal(v any) (map[string]string, error) {
	ref := &Reflect{Type: reflect.TypeOf(v)}
	if ref.Type.Kind() == reflect.Pointer {
		ref.Type = ref.Type.Elem()
	}
	if ref.Type == nil || ref.Type.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	return ref.encode(v)
}
