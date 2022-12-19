package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/liasece/gocoder"
)

type ASTPkg struct {
	name string
	node ast.Node
}

type ASTCoder struct {
	fset         *token.FileSet
	pkgs         []*ASTPkg
	importPkgs   map[string]string
	DecodedTypes map[string]*ASTTyped
}

type ASTTyped struct {
	gocoder.Type
}

func NewASTCoder(paths ...string) (*ASTCoder, error) {
	fset := token.NewFileSet()
	ps := make([]*ASTPkg, 0)
	for _, path := range paths {
		pathSS := strings.Split(path, ",")
		for _, path := range pathSS {
			if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
				// from path
				pkgs, err := ParseDir(fset, path, nil, parser.AllErrors)
				if err != nil {
					return nil, err
				}
				for k, v := range pkgs {
					ps = append(ps, &ASTPkg{
						name: k,
						node: v,
					})
				}
			} else {
				node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
				if err != nil {
					return nil, err
				}
				ps = append(ps, &ASTPkg{
					name: node.Name.Name,
					node: node,
				})
			}
		}
	}
	return &ASTCoder{
		fset:         fset,
		pkgs:         ps,
		importPkgs:   make(map[string]string),
		DecodedTypes: make(map[string]*ASTTyped),
	}, nil
}

func GetTypeFromSource(path string, typeName string, opts ...*gocoder.ToCodeOption) (gocoder.Type, error) {
	opt := gocoder.MergeToCodeOpt(opts...)
	c, err := NewASTCoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetType(typeName, opt)
}

func GetInterfaceFromSource(path string, typeName string, opts ...*gocoder.ToCodeOption) (gocoder.Interface, error) {
	opt := gocoder.MergeToCodeOpt(opts...)
	c, err := NewASTCoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetInterface(typeName, opt)
}

func GetMethodsFromSource(path string, typeName string, opts ...*gocoder.ToCodeOption) ([]gocoder.Func, error) {
	opt := gocoder.MergeToCodeOpt(opts...)
	c, err := NewASTCoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetMethods(typeName, opt)
}
