package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/liasece/gocoder"
)

func ParseDir(fset *token.FileSet, path string, filter func(fs.FileInfo) bool, mode parser.Mode) (pkgs map[string]*ast.Package, first error) {
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	pkgs = make(map[string]*ast.Package)
	for _, d := range list {
		if d.IsDir() {
			ps, err := ParseDir(fset, filepath.Join(path, d.Name()), filter, mode)
			if err != nil {
				return nil, err
			}
			for k, v := range ps {
				pkgs[k] = v
			}
			continue
		}
		if !strings.HasSuffix(d.Name(), ".go") {
			continue
		}
		if filter != nil {
			info, err := d.Info()
			if err != nil {
				return nil, err
			}
			if !filter(info) {
				continue
			}
		}
		filename := filepath.Join(path, d.Name())
		if src, err := parser.ParseFile(fset, filename, nil, mode); err == nil {
			name := src.Name.Name
			pkg, found := pkgs[name]
			if !found {
				pkg = &ast.Package{
					Name:  name,
					Files: make(map[string]*ast.File),
				}
				pkgs[name] = pkg
			}
			pkg.Files[filename] = src
		} else if first == nil {
			first = err
		}
	}

	return
}

// wrap a function to fulfill ast.Visitor interface
type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

func TypeStringToReflectKind(str string) reflect.Kind {
	switch str {
	case "string":
		return reflect.String
	}
	return reflect.Invalid
}

func toTypeStr(pkg string, name string) string {
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

func pkgInReference(str string) string {
	ss := strings.Split(str, "/")
	return ss[len(ss)-1]
}

var InterfaceType reflect.Type

func init() {
	type T struct {
		A interface{}
	}
	InterfaceType = reflect.ValueOf(T{}).Field(0).Type()
}

func TypeStringToZeroInterface(str string) gocoder.Type {
	switch str {
	case "bool":
		return gocoder.MustToType(false)
	case "byte":
		return gocoder.MustToType(byte(0))
	case "int":
		return gocoder.MustToType(int(0))
	case "int8":
		return gocoder.MustToType(int8(0))
	case "int16":
		return gocoder.MustToType(int16(0))
	case "int32":
		return gocoder.MustToType(int32(0))
	case "int64":
		return gocoder.MustToType(int64(0))
	case "uint":
		return gocoder.MustToType(uint(0))
	case "uint8":
		return gocoder.MustToType(uint8(0))
	case "uint16":
		return gocoder.MustToType(uint16(0))
	case "uint32":
		return gocoder.MustToType(uint32(0))
	case "uint64":
		return gocoder.MustToType(uint64(0))
	case "float32":
		return gocoder.MustToType(float32(0))
	case "float64":
		return gocoder.MustToType(float64(0))
	case "complex64":
		return gocoder.MustToType(complex64(0))
	case "complex128":
		return gocoder.MustToType(complex128(0))
	case "string":
		return gocoder.MustToType("")
	case "time.Time":
		return gocoder.MustToType(time.Time{})

	case "[]bool":
		return gocoder.MustToType([]bool{})
	case "[]int":
		return gocoder.MustToType([]int{})
	case "[]int8":
		return gocoder.MustToType([]int8{})
	case "[]int16":
		return gocoder.MustToType([]int16{})
	case "[]int32":
		return gocoder.MustToType([]int32{})
	case "[]int64":
		return gocoder.MustToType([]int64{})
	case "[]uint":
		return gocoder.MustToType([]uint{})
	case "[]uint8":
		return gocoder.MustToType([]uint8{})
	case "[]uint16":
		return gocoder.MustToType([]uint16{})
	case "[]uint32":
		return gocoder.MustToType([]uint32{})
	case "[]uint64":
		return gocoder.MustToType([]uint64{})
	case "[]float32":
		return gocoder.MustToType([]float32{})
	case "[]float64":
		return gocoder.MustToType([]float64{})
	case "[]complex64":
		return gocoder.MustToType([]complex64{})
	case "[]complex128":
		return gocoder.MustToType([]complex128{})
	case "[]string":
		return gocoder.MustToType([]string{})
	case "[]time.Time":
		return gocoder.MustToType([]time.Time{})

	case "interface{}":
		return gocoder.MustToType(InterfaceType)

	case "error":
		return gocoder.NewTypeName("error")

	case "context.Context":
		return gocoder.NewTypeDetail("context", "Context")
	}
	return nil
}
