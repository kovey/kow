package json

import "encoding/json"

var Marshal func(v any) ([]byte, error) = json.Marshal
var Unmarshal func(data []byte, v any) error = json.Unmarshal

func Init(m func(v any) ([]byte, error), u func(data []byte, v any) error) {
	Marshal = m
	Unmarshal = u
}
