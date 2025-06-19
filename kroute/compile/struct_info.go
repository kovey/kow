package compile

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type structInfo struct {
	tpl  *templateKroute
	file *ast.File
	has  bool
}

func newStructInfo(file *ast.File) *structInfo {
	return &structInfo{tpl: &templateKroute{}, file: file}
}

func (s *structInfo) replace() error {
	buf, err := s.tpl.Parse()
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", buf, 0)
	if err != nil {
		return err
	}

	for _, d := range file.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok {
			s.file.Decls = append(s.file.Decls, fn)
			continue
		}

	}

	ims := s.file.Decls[0].(*ast.GenDecl)
	for _, im := range file.Imports {
		has := false
		for _, imed := range s.file.Imports {
			if imed.Path.Value == im.Path.Value {
				has = true
				break
			}
		}

		if !has {
			ims.Specs = append(ims.Specs, im)
		}
	}

	return nil
}

func isKrouteStruct(gn *ast.GenDecl, flag *string) bool {
	for i := len(gn.Doc.List) - 1; i >= 0; i-- {
		if strings.HasPrefix(strings.Trim(gn.Doc.List[i].Text, " "), tag_kroute) {
			return true
		}
	}

	return false
}

func _parseStruct(pkg *_packageFile) *_structInfos {
	sts := newStructInfos()
	for file, f := range pkg.files {
		var _struct = newStructInfo(f)
		for _, d := range f.Decls {
			if gn, ok := d.(*ast.GenDecl); ok && gn.Doc != nil {
				var docs = make([]string, len(gn.Doc.List))
				has := false
				for i, doc := range gn.Doc.List {
					docs[i] = doc.Text
					if strings.HasPrefix(strings.Trim(gn.Doc.List[i].Text, " "), tag_kroute) {
						has = true
					}
				}

				if !has {
					continue
				}

				at := &Ast{}
				if err := at.Parse(docs); err != nil {
					continue
				}
				_struct.tpl.Routers = append(_struct.tpl.Routers, at.router(_struct.tpl))
				_struct.has = true
			}
		}

		if !_struct.has {
			continue
		}

		sts.files[file] = f
		sts.structs[file] = _struct
	}

	return sts
}

type _structInfos struct {
	structs    map[string]*structInfo
	files      map[string]*ast.File
	printerCfg *printer.Config
}

func newStructInfos() *_structInfos {
	return &_structInfos{structs: make(map[string]*structInfo), files: make(map[string]*ast.File), printerCfg: &printer.Config{Tabwidth: 8, Mode: printer.SourcePos}}
}

func (s *_structInfos) replace(fset *token.FileSet, tempDir string, args []string) error {
	for _, st := range s.structs {
		if err := st.replace(); err != nil {
			return err
		}
	}

	for file, f := range s.files {
		originPath := file
		tgDir := path.Join(tempDir, os.Getenv("TOOLEXEC_IMPORTPATH"))
		buffer := bytes.NewBuffer(nil)
		if err := s.printerCfg.Fprint(buffer, fset, f); err != nil {
			return err
		}

		_ = os.MkdirAll(tgDir, 0777)
		tmpEntryFile := path.Join(tgDir, filepath.Base(originPath))
		if err := os.WriteFile(tmpEntryFile, buffer.Bytes(), 0777); err != nil {
			return err
		}

		for i := range args {
			if args[i] == originPath {
				args[i] = tmpEntryFile
				break
			}
		}
	}

	return nil
}
