package cde

import (
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/ast"
)

func TypeStringToZeroInterface(str string) gocoder.Type {
	return ast.TypeStringToZeroInterface(str)
}

func GetTypeFromSource(path string, typeName string, opts ...*gocoder.ToCodeOption) (gocoder.Type, error) {
	opt := gocoder.MergeToCodeOpt(opts...)
	c, err := ast.NewASTCoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetType(typeName, opt)
}

func GetInterfaceFromSource(path string, typeName string, opts ...*gocoder.ToCodeOption) (gocoder.Interface, error) {
	opt := gocoder.MergeToCodeOpt(opts...)
	c, err := ast.NewASTCoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetInterface(typeName, opt)
}

func GetMethodsFromSource(path string, typeName string, opts ...*gocoder.ToCodeOption) ([]gocoder.Func, error) {
	opt := gocoder.MergeToCodeOpt(opts...)
	c, err := ast.NewASTCoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetMethods(typeName, opt)
}
