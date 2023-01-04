package cde

import (
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/ast"
)

func TypeStringToZeroInterface(str string) gocoder.Type {
	return ast.TypeStringToZeroInterface(str)
}

func GetTypeFromSource(path string, typeName string) (gocoder.Type, error) {
	c, err := ast.NewCodeDecoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetType(typeName), nil
}

func GetInterfaceFromSource(path string, typeName string) (gocoder.Type, error) {
	c, err := ast.NewCodeDecoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetInterface(typeName), nil
}

func GetMethodsFromSource(path string, typeName string) ([]gocoder.Func, error) {
	c, err := ast.NewCodeDecoder(path)
	if err != nil {
		return nil, err
	}
	return c.GetMethods(typeName), nil
}
