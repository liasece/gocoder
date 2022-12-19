package ast

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// get target go file's full package name, like "github.com/liasece/gocoder/ast" "ast".
// is the package need rename, like "github.com/liasece/gocoder/ast" "ast_test".
// will lock up the go.mod and join the module name. The filePath must be a go file.
func GetGoFileFullPackage(filePath string) (pkg string, alias string) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, file, parser.PackageClauseOnly)
	if err != nil {
		return "", ""
	}

	moduleName := GetGoModModuleName(filepath.Dir(filePath))
	if moduleName == "" {
		return f.Name.Name, ""
	}
	return moduleName, f.Name.Name
}

func getGoModModuleName(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path does not exist
		return "", err
	}
	file, err := os.Open(filepath.Join(path, "go.mod"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			mod, err := getGoModModuleName(filepath.Dir(path))
			if err != nil {
				return "", err
			}
			return mod + "/" + filepath.Base(path), nil
		}
		return "", err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	module := modfile.ModulePath(bytes)
	return module, nil
}

func GetGoModModuleName(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	moduleName, err := getGoModModuleName(absPath)
	if err != nil {
		return ""
	}
	return moduleName
}

func Parse(fset *token.FileSet, path string, filter func(fs.FileInfo) bool, mode parser.Mode) (pkgs *Packages, first error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	pkgs = &Packages{}
	if info.IsDir() {
		{
			// this package
			pkg, err := parser.ParseDir(fset, path, filter, mode)
			if err != nil {
				return nil, err
			}
			for _, v := range pkg {
				pkg, alias := "", ""
				for k := range v.Files {
					pkg, alias = GetGoFileFullPackage(k)
					break
				}
				if pkg != "" && alias != "" {
					pkgs.Add(&Package{
						Name:    pkg,
						Alias:   alias,
						Package: v,
					})
				}
			}
		}
		list, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, d := range list {
			if d.IsDir() {
				ps, err := Parse(fset, filepath.Join(path, d.Name()), filter, mode)
				if err != nil {
					return nil, err
				}
				pkgs.Add(ps.List...)
				continue
			}
		}
	} else {
		// this file
		if src, err := parser.ParseFile(fset, path, nil, mode); err == nil {
			name := src.Name.Name
			pkg := &ast.Package{
				Name: name,
				Files: map[string]*ast.File{
					path: src,
				},
			}
			pkgStr, alias := GetGoFileFullPackage(path)
			if pkgStr != "" && alias != "" {
				pkgs.Add(&Package{
					Name:    pkgStr,
					Alias:   alias,
					Package: pkg,
				})
			}
		}
	}

	return
}
