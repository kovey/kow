package compile

import (
	"go/format"
	"html/template"
	"strings"
)

const (
	import_kow      = `"github.com/kovey/kow"`
	template_kroute = `package compile
import(
	"github.com/kovey/kow"
	"github.com/kovey/kow/router"
	{{range .Imports}}{{"\r\n"}}"{{.}}"{{end}}
)

func init() {
	var r = router.RouterInterface
{{range .Routers}}
	r = kow.{{.Method | safe}}("{{.Path | safe}}", {{.Constructor | safe}}).Data({{.ReqData | safe}})
	{{range .Rules}}{{"\r\n"}}r.Rule("{{.Name | safe}}"{{range .Args}}{{","}}"{{. | safe}}"{{end}}){{end}}
	{{range .Middlewares}}{{"\r\n"}}r.Middleware({{. | safe}}){{end}}
	{{"\r\n"}}
{{end}}
}
`
)

const (
	tag_kroute     = "//go:kroute " // go:kroute router("PUT", "/path", newAction(), &ReqData{}), rule("name", ["minlen:int:1", "maxlen:int:127"]), middleware(newAuth(), middlewares.NewLog()), import("path1", "path2")
	tag_router     = "router"       // router("PUT", "/path", newAction(), &ReqData{})
	tag_rule       = "rule"         // rule("name", ["minlen:int:1", "maxlen:int:127"])
	tag_middleware = "middleware"   // middleware(newAuth(), middlewares.NewLog())
	tag_import     = "import"       // import("path1", "path2")
)

type router struct {
	Method      string
	Path        string
	Constructor string
	ReqData     string
	Rules       []*rule
	Middlewares []string
}

type rule struct {
	Name string
	Args []string
}

type _column struct {
	name string
	tag  string
}

type columnInfo struct {
	Name string
	Type string
}

type templateKroute struct {
	Routers []*router
	Imports []string
}

func (tk *templateKroute) addImport(im string) {
	has := false
	for _, i := range tk.Imports {
		if i == im {
			has = true
			break
		}
	}

	if !has {
		tk.Imports = append(tk.Imports, im)
	}
}

func (tk *templateKroute) Parse() ([]byte, error) {
	t := template.Must(template.New("main_tpl").Funcs(template.FuncMap{"safe": func(tag string) template.HTML {
		return template.HTML(tag)
	}}).Parse(template_kroute))
	builder := strings.Builder{}
	if err := t.Execute(&builder, tk); err != nil {
		return nil, err
	}

	return format.Source([]byte(builder.String()))
}
