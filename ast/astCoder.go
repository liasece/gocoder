package ast

import (
	"go/parser"
	"go/token"
	"strings"

	"github.com/liasece/gocoder"
)

type CodeDecoder struct {
	fset         *token.FileSet
	pkgs         *Packages
	DecodedTypes map[string]*LoadedType
}

type LoadedType struct {
	gocoder.Type
}

func NewCodeDecoder(paths ...string) (*CodeDecoder, error) {
	fset := token.NewFileSet()
	ps := &Packages{
		List: nil,
	}
	for _, path := range paths {
		pathSS := strings.Split(path, ",")
		for _, path := range pathSS {
			// from path
			pkgs, err := Parse(fset, path, nil, parser.AllErrors|parser.ParseComments)
			if err != nil {
				return nil, err
			}
			ps.MergeFrom(pkgs)
		}
	}
	return &CodeDecoder{
		fset:         fset,
		pkgs:         ps,
		DecodedTypes: make(map[string]*LoadedType),
	}, nil
}

func GetTypeFromSource(path string, typeName string) (gocoder.Type, error) {
	c, err := NewCodeDecoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetType(typeName), nil
}

func GetInterfaceFromSource(path string, typeName string) (gocoder.Type, error) {
	c, err := NewCodeDecoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetInterface(typeName), nil
}

func GetMethodsFromSource(path string, typeName string) ([]gocoder.Func, error) {
	c, err := NewCodeDecoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetMethods(typeName), nil
}
