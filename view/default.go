package view

import (
	"fmt"
	"html/template"
	"io"
)

type Default struct {
	data any
	tpl  *template.Template
}

func NewDefault(tpl *template.Template) *Default {
	return &Default{tpl: tpl}
}

func (d *Default) Data(data any) {
	d.data = data
}

func (d *Default) Load(path string) error {
	tpl, err := template.ParseGlob(path)
	if err != nil {
		return err
	}

	d.tpl = tpl
	return nil
}

func (d *Default) Parse(writer io.Writer) error {
	if d.tpl == nil {
		return fmt.Errorf("template is nil")
	}

	return d.tpl.Execute(writer, d.data)
}
