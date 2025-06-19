package ast

import (
	"fmt"
	"strings"
)

// router("PUT", "/path", newAction(), &ReqData{})
type Router struct {
	Begin       byte
	Method      *Arg
	Path        *Arg
	End         byte
	Constructor *Expr
	ReqData     *Expr
}

func (r *Router) Completed() bool {
	return r.Begin > 0 && r.End > 0
}

func (r *Router) Parse(b byte) error {
	switch b {
	case '(':
		if r.Begin > 0 {
			if r.Constructor == nil {
				return fmt.Errorf("parse router unexpect (")
			}

			if !r.Constructor.check() {
				return r.Constructor.Parse(b)
			}

			if r.ReqData == nil {
				return fmt.Errorf("parse router unexpect (")
			}

			if !r.ReqData.check() {
				return r.ReqData.Parse(b)
			}
		}

		r.Begin = b
		r.Method = &Arg{}
	case ')':
		if r.Begin == 0 {
			return fmt.Errorf("parse router unexpect )")
		}

		if r.Constructor != nil && !r.Constructor.check() {
			return r.Constructor.Parse(b)
		}

		if r.ReqData != nil && !r.ReqData.check() {
			return r.ReqData.Parse(b)
		}

		r.End = b
	case '"':
		if r.Method == nil {
			return fmt.Errorf("unexpect \"")
		}

		if !r.Method.check() {
			return r.Method.parse(b)
		}

		if r.Path == nil {
			return fmt.Errorf("unexpect \"")
		}

		if !r.Path.check() {
			return r.Path.parse(b)
		}

		return fmt.Errorf("parse router failure")
	case ',':
		if r.Method == nil {
			return fmt.Errorf(`parse router unexpect ","`)
		}

		if r.Path == nil {
			r.Path = &Arg{}
			return nil
		}

		if r.Constructor == nil {
			r.Constructor = &Expr{}
			return nil
		}

		if r.ReqData == nil {
			r.ReqData = &Expr{}
			return nil
		}
		return fmt.Errorf(`parse router unexpect ","`)
	case ' ':
		if r.Method == nil {
			return nil
		}

		if !r.Method.check() {
			return r.Method.parse(b)
		}

		if r.Path == nil {
			return nil
		}

		if !r.Path.check() {
			return r.Path.parse(b)
		}

		if r.Constructor == nil {
			return nil
		}

		if !r.Constructor.check() {
			return r.Constructor.Parse(b)
		}

		if r.ReqData == nil {
			return nil
		}

		if !r.ReqData.check() {
			return r.ReqData.Parse(b)
		}
	default:
		if r.Method == nil {
			return fmt.Errorf(`parse router unexpect "%c"`, b)
		}
		if !r.Method.check() {
			return r.Method.parse(b)
		}

		if r.Path == nil {
			return fmt.Errorf(`parse router unexpect "%c"`, b)
		}

		if !r.Path.check() {
			return r.Path.parse(b)
		}

		if r.Constructor == nil {
			return fmt.Errorf(`parse router unexpect "%c"`, b)
		}

		if !r.Constructor.check() {
			return r.Constructor.Parse(b)
		}
		if r.ReqData == nil {
			return fmt.Errorf(`parse router unexpect "%c"`, b)
		}

		if !r.ReqData.check() {
			return r.ReqData.Parse(b)
		}

		return fmt.Errorf(`parse router unexpect "%c"`, b)
	}

	return nil
}

type Expr struct {
	Star     byte
	Name     string
	Left     byte
	Right    byte
	vBuidler strings.Builder
}

func (e *Expr) String() string {
	e.vBuidler.Reset()
	if e.Star > 0 {
		e.vBuidler.WriteByte(e.Star)
	}
	e.vBuidler.WriteString(e.Name)
	e.vBuidler.WriteByte(e.Left)
	e.vBuidler.WriteByte(e.Right)
	return e.vBuidler.String()
}

func (e *Expr) isNil() bool {
	return e.Name == "nil"
}

func (e *Expr) check() bool {
	return e.Left > 0 && e.Right > 0
}

func (e *Expr) Parse(b byte) error {
	switch b {
	case '&':
		if e.Star > 0 || e.Left > 0 || e.vBuidler.String() != "" {
			return fmt.Errorf("expr unexpect &")
		}

		e.Star = b
	case '(', '{':
		if e.Left > 0 {
			return fmt.Errorf("expr unexpect %c", b)
		}

		e.Left = b
		e.Name = e.vBuidler.String()
		if e.Name == "" {
			return fmt.Errorf("expr unexpect %c", b)
		}
	case ')', '}':
		if e.Right > 0 || e.Left == 0 {
			return fmt.Errorf("expr unexpect %c", b)
		}

		switch b {
		case ')':
			if e.Left != '(' {
				return fmt.Errorf("expr unexpect %c", b)
			}
		case '}':
			if e.Left != '{' {
				return fmt.Errorf("expr unexpect %c", b)
			}
		}

		e.Right = b
	case ' ':
		return nil
	default:
		if e.Left > 0 || e.Name != "" {
			return fmt.Errorf("expr unexpect %c", b)
		}

		return e.vBuidler.WriteByte(b)
	}

	return nil
}
