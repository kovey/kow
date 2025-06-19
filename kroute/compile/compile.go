package compile

import (
	"fmt"
	"go/token"
	"path/filepath"
	"strings"
)

func Compile(projectDir, tempDir string, args []string) error {
	packageInfo, err := parsePackageInfo(projectDir, "")
	if err != nil {
		return err
	}
	if packageInfo.Module.Path == "" {
		return fmt.Errorf("module in %s not found", projectDir)
	}
	files := make([]string, 0, len(args))
	projectName := packageInfo.Module.Path
	packageName := packageInfo.Name
	for i, arg := range args {
		if arg == "-p" && i+1 < len(args) {
			packageName = args[i+1]
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		if strings.HasPrefix(arg, projectDir+string(filepath.Separator)) && strings.HasSuffix(arg, ".go") {
			files = args[i:]
			break
		}
	}

	files = args
	if (packageName != "main" && !strings.HasPrefix(packageName, projectName)) || len(files) == 0 {
		return nil
	}

	fset := token.NewFileSet()
	pkg := newPackageFile(projectDir, packageInfo.Name)
	if err := pkg.ParseFile(fset, files); err != nil {
		return err
	}

	sts := _parseStruct(pkg)
	return sts.replace(fset, tempDir, args)
}
