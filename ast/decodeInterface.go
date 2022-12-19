package ast

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func (c *ASTCoder) GetInterfaceFromASTInterfaceType(name string, st *ast.InterfaceType, opt *gocoder.ToCodeOption) (gocoder.Interface, error) {
	fs, err := c.GetFuncsFromASTFieldList(nil, st.Methods, opt)
	if err != nil {
		return nil, err
	}
	res := gocoder.NewInterface(name, fs)
	return res, nil
}

func (c *ASTCoder) GetInterface(name string, opt *gocoder.ToCodeOption) (gocoder.Interface, error) {
	var resType gocoder.Interface
	var resErr error
	typeTypeName := name
	// typePkgName := ""
	if ss := strings.Split(name, "."); len(ss) == 2 {
		// typePkgName = ss[0]
		typeTypeName = ss[1]
	}
	for _, pkgV := range c.pkgs {
		_, node := pkgV.name, pkgV.node
		ast.Walk((walker)(func(node ast.Node) bool {
			if node == nil {
				return true
			}
			{
				// add import
				if ts, ok := node.(*ast.ImportSpec); ok {
					path := strings.ReplaceAll(ts.Path.Value, "\"", "")
					ss := strings.Split(path, "/")
					if len(ss) > 0 {
						c.importPkgs[ss[len(ss)-1]] = path
					}
				}
			}
			ts, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if ts.Name.Name == typeTypeName {
				// found target type
				if st, ok := ts.Type.(*ast.InterfaceType); ok {
					resType, resErr = c.GetInterfaceFromASTInterfaceType(ts.Name.Name, st, opt)
					return false
				}
				log.Error("GetInterface ts.Name.Name == name but type unknown", log.Any("type", reflect.TypeOf(ts.Type)))
			}
			return true
		}), node)
		if resErr != nil || resType != nil {
			break
		}
	}
	if resType == nil {
		log.Warn("GetInterface not found type", log.Any("name", name), log.Any("typeTypeName", typeTypeName))
	}
	return resType, resErr
}
