package ast

import (
	"fmt"
	"strings"
)

// rule("Name", ["minlen:int:1", "maxlen:int:127"])
type Rule struct {
	Begin byte
	Name  *Arg
	Args  *Args
	End   byte
}

func (r *Rule) Completed() bool {
	return r.Begin > 0 && r.End > 0
}

func (r *Rule) Parse(b byte) error {
	switch b {
	case '(':
		if r.Begin > 0 {
			return fmt.Errorf("parse rule unexpect (")
		}

		r.Begin = b
		r.Name = &Arg{}
	case ')':
		if r.Begin == 0 {
			return fmt.Errorf("parse rule unexpect )")
		}

		r.End = b
	case '[':
		if r.Name == nil || r.Args == nil {
			return fmt.Errorf("parse rule unexpect [")
		}
		return r.Args.parse(b)
	case ']':
		if r.Args == nil {
			return fmt.Errorf("parse rule unexpect ]")
		}
		if err := r.Args.parse(b); err != nil {
			return err
		}
		if !r.Args.check() {
			return fmt.Errorf("parse rule unexpect ]")
		}
	case '"':
		if r.Name == nil {
			return fmt.Errorf("unexpect \"")
		}

		if !r.Name.check() {
			return r.Name.parse(b)
		}

		if r.Args == nil {
			return fmt.Errorf("parse rule unexpect %c", b)
		}

		return r.Args.parse(b)
	case ',':
		if r.Name == nil || !r.Name.check() {
			return fmt.Errorf(`parse rule unexpect ","`)
		}

		if r.Args == nil {
			r.Args = &Args{}
			return nil
		}

		return r.Args.parse(b)
	case ' ':
		if r.Begin == 0 {
			return nil
		}

		if !r.Name.check() {
			return fmt.Errorf(`unexpect " "`)
		}

		if r.Args == nil {
			return nil
		}

		return r.Args.parse(b)
	default:
		if r.Name == nil {
			return fmt.Errorf("unexpect %c", b)
		}

		if !r.Name.check() {
			return r.Name.parse(b)
		}

		if r.Args != nil {
			return r.Args.parse(b)
		}
	}

	return nil
}

type Args struct {
	Begin   byte
	Args    []*Arg
	End     byte
	current *Arg
}

func (f *Args) check() bool {
	return f.Begin > 0 && f.End > 0 && f.current == nil
}

func (f *Args) parse(b byte) error {
	switch b {
	case '[':
		if f.Begin != 0 {
			return fmt.Errorf("unexpect [")
		}
		f.Begin = b
		f.current = &Arg{}
	case ']':
		f.End = b
		if f.Begin == 0 {
			return fmt.Errorf("Rule Args Begin not found")
		}
	case ',':
		if f.Begin == 0 {
			return fmt.Errorf("unexpect ,")
		}

		if f.current != nil {
			return fmt.Errorf("Rule arg error")
		}

		f.current = &Arg{}
	case ' ':
		if f.current == nil {
			return nil
		}
	default:
		if f.current == nil {
			return fmt.Errorf("unexpect %c", b)
		}

		if err := f.current.parse(b); err != nil {
			return err
		}

		if f.current.check() {
			f.Args = append(f.Args, f.current)
			f.current = nil
		}
	}

	return nil
}

type Arg struct {
	Begin    byte
	Value    string
	End      byte
	vBuilder strings.Builder
}

func (f *Arg) check() bool {
	return f.Begin > 0 && f.End > 0
}

func (f *Arg) parse(b byte) error {
	switch b {
	case ' ':
		if f.Begin > 0 {
			return fmt.Errorf("arg format error")
		}

		return nil
	case '"':
		if f.Begin == 0 {
			f.Begin = b
			return nil
		}

		f.Value = f.vBuilder.String()
		f.End = b
	default:
		if f.check() {
			return nil
		}

		if f.Begin == 0 {
			return fmt.Errorf("Begin not found")
		}

		return f.vBuilder.WriteByte(b)
	}

	return nil
}
