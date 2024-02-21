package xml

import "encoding/xml"

var Marshal func(v any) ([]byte, error) = xml.Marshal
var Unmarshal func(data []byte, v any) error = xml.Unmarshal

func Init(m func(v any) ([]byte, error), u func(data []byte, v any) error) {
	Marshal = m
	Unmarshal = u
}
