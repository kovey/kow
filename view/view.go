package view

import "io"

type ViewInterface interface {
	Load(path string) error
	Data(data any)
	Parse(writer io.Writer) error
}
