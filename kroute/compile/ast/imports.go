package ast

import "fmt"

// import("path1", "path2")
type Import struct {
	Args       []*Arg
	Begin      byte
	End        byte
	currentArg *Arg
}

func (i *Import) check() bool {
	return i.Begin > 0 && i.End > 0
}

func (i *Import) Completed() bool {
	return i.check()
}

func (i *Import) Parse(b byte) error {
	switch b {
	case '(':
		if i.Begin > 0 {
			return fmt.Errorf("unexpect (")
		}

		i.Begin = b
		i.currentArg = &Arg{}
	case ')':
		if i.End > 0 {
			return fmt.Errorf("unexpect )")
		}

		i.End = b
	case ',':
		if i.currentArg != nil || i.Begin == 0 {
			return fmt.Errorf("unexpect ,")
		}

		i.currentArg = &Arg{}
	case ' ':
		if i.currentArg == nil {
			return nil
		}
		return i.currentArg.parse(b)
	default:
		if i.currentArg == nil {
			return fmt.Errorf("unexpect %c", b)
		}
		if err := i.currentArg.parse(b); err != nil {
			return err
		}

		if i.currentArg.check() {
			i.Args = append(i.Args, i.currentArg)
			i.currentArg = nil
		}
	}

	return nil
}
