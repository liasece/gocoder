package ast

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func (c *CodeDecoder) GetInterfaceFromASTInterfaceType(ctx DecoderContext, st *ast.InterfaceType) gocoder.Type {
	fs := c.GetFuncsFromASTFieldList(ctx, nil, st.Methods)
	res := gocoder.NewInterface(ctx.GetBuildingItemName(), fs)
	return res
}

func (c *CodeDecoder) GetInterface(name string) gocoder.Type {
	var resType gocoder.Type
	typeTypeName := name
	typePkg := ""
	if index := strings.LastIndex(name, "."); index > 0 && index < len(name)-1 {
		typePkg = name[:index]
		typeTypeName = name[index+1:]
	}
	for _, pkgV := range c.pkgs.List {
		if typePkg != "" {
			if typePkg != pkgV.Name && typePkg != pkgV.Alias {
				continue
			}
		}
		for _, astFile := range pkgV.Package.Files {
			ast.Walk((walker)(func(node ast.Node) bool {
				if node == nil {
					return true
				}
				ts, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}
				if ts.Name.Name == typeTypeName {
					// found target type
					if st, ok := ts.Type.(*ast.InterfaceType); ok {
						ctx := NewDecoderContextByAstFile(pkgV.Name, typeTypeName, astFile)
						resType = c.GetInterfaceFromASTInterfaceType(ctx, st)
						if resType != nil {
							return false
						}
					} else {
						log.Error("GetInterface ts.Name.Name == name but type unknown", log.Any("type", reflect.TypeOf(ts.Type)))
					}
				}
				return true
			}), astFile)
			if resType != nil {
				break
			}
		}
		if resType != nil {
			break
		}
	}
	if resType == nil {
		log.Warn("GetInterface not found type", log.Any("name", name), log.Any("typeTypeName", typeTypeName))
	}
	return resType
}
