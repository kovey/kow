package ast

import "fmt"

// middleware(newAuth(), middlewares.NewLog())
type Middleware struct {
	Args    []*Expr
	Begin   byte
	End     byte
	current *Expr
}

func (m *Middleware) check() bool {
	return m.Begin > 0 && m.End > 0
}

func (m *Middleware) Completed() bool {
	return m.check()
}

func (m *Middleware) Parse(b byte) error {
	switch b {
	case '(':
		if m.Begin > 0 {
			if m.current == nil || m.current.Left > 0 {
				return fmt.Errorf("middleware unexpect (")
			}

			return m.current.Parse(b)
		}

		m.Begin = b
		m.current = &Expr{}
	case ')':
		if m.current != nil {
			if err := m.current.Parse(b); err != nil {
				return err
			}

			if m.current.check() {
				m.Args = append(m.Args, m.current)
				m.current = nil
				return nil
			}
		}

		if m.End > 0 {
			return fmt.Errorf("middleware unexpect )")
		}

		m.End = b
	case ',':
		if m.current != nil {
			return fmt.Errorf("middleware unexpect ,")
		}

		m.current = &Expr{}
	case ' ':
		if m.current == nil {
			return nil
		}

		return m.current.Parse(b)
	default:
		if m.current == nil {
			return fmt.Errorf("middleware unexpect %c", b)
		}

		if err := m.current.Parse(b); err != nil {
			return err
		}

		if m.current.check() {
			m.Args = append(m.Args, m.current)
			m.current = nil
		}
	}

	return nil
}
