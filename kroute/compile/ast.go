package compile

import (
	"fmt"
	"strings"

	"github.com/kovey/kow/kroute/compile/ast"
)

type Ast struct {
	Router      *ast.Router
	Rules       []*ast.Rule
	Middlewares []*ast.Middleware
	Imports     []*ast.Import
}

func (a *Ast) router(tk *templateKroute) *router {
	t := &router{Method: a.Router.Method.Value, Path: a.Router.Path.Value, Constructor: a.Router.Constructor.String(), ReqData: a.Router.ReqData.String()}
	for _, r := range a.Rules {
		rr := &rule{Name: r.Name.Value}
		for _, arg := range r.Args.Args {
			rr.Args = append(rr.Args, arg.Value)
		}

		t.Rules = append(t.Rules, rr)
	}

	for _, m := range a.Middlewares {
		for _, ex := range m.Args {
			t.Middlewares = append(t.Middlewares, ex.String())
		}
	}

	for _, i := range a.Imports {
		for _, ex := range i.Args {
			tk.addImport(ex.Value)
		}
	}
	return t
}

func (a *Ast) Parse(docments []string) error {
	for _, doc := range docments {
		if !strings.HasPrefix(doc, tag_kroute) {
			continue
		}

		doc = strings.ReplaceAll(doc, tag_kroute, "")
		var builder strings.Builder
		buff := []byte(doc)
		count := len(buff)
		var i = 0
	parse:
		builder.Reset()
		for i < count {
			if buff[i] == ' ' {
				if builder.Len() > 0 {
					builder.WriteByte(buff[i])
				}
			} else {
				builder.WriteByte(buff[i])
			}
			i++
			switch builder.String() {
			case tag_import, tag_middleware, tag_rule, tag_router:
				if err := a.parse(&i, buff, builder.String(), count); err != nil {
					return err
				}
				goto parse
			}
		}
	}
	return nil
}

func (a *Ast) parse(start *int, buff []byte, tag string, end int) error {
	switch tag {
	case tag_import:
		im := &ast.Import{}
		for *start < end {
			if err := im.Parse(buff[*start]); err != nil {
				return err
			}
			*start++
			if im.Completed() {
				a.Imports = append(a.Imports, im)
				return nil
			}
		}
	case tag_router:
		if a.Router != nil {
			return fmt.Errorf("router is repeated")
		}

		r := &ast.Router{}
		for *start < end {
			if err := r.Parse(buff[*start]); err != nil {
				return err
			}
			*start++
			if r.Completed() {
				a.Router = r
				return nil
			}
		}
	case tag_rule:
		r := &ast.Rule{}
		for *start < end {
			if err := r.Parse(buff[*start]); err != nil {
				return err
			}
			*start++
			if r.Completed() {
				a.Rules = append(a.Rules, r)
				return nil
			}
		}
	case tag_middleware:
		m := &ast.Middleware{}
		for *start < end {
			if err := m.Parse(buff[*start]); err != nil {
				return err
			}
			*start++
			if m.Completed() {
				a.Middlewares = append(a.Middlewares, m)
				return nil
			}
		}
	}
	return nil
}
