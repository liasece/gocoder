package ast

import (
	"go/ast"
	"strings"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func (c *CodeDecoder) GetMethods(reviverTypeName string) []gocoder.Func {
	var res []gocoder.Func
	typeTypeName := reviverTypeName
	typePkg := ""
	if index := strings.LastIndex(reviverTypeName, "."); index > 0 && index < len(reviverTypeName)-1 {
		typePkg = reviverTypeName[:index]
		typeTypeName = reviverTypeName[index+1:]
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
				ts, ok := node.(*ast.FuncDecl)
				if !ok {
					return true
				}
				if ts.Recv == nil || len(ts.Recv.List) == 0 {
					return true
				}
				ctx := NewDecoderContextByAstFile(pkgV.Name, typeTypeName, astFile)
				reviverType := c.GetTypeFromASTNode(ctx, ts.Recv.List[0].Type)
				name := ""
				if reviverType != nil {
					name = reviverType.String()
				}
				if name != reviverTypeName && name != "*"+reviverTypeName {
					return true
				}
				fn := c.GetFuncsFromASTFuncDecl(ctx, ts)
				if fn != nil {
					res = append(res, fn)
				}
				return true
			}), astFile)
		}
	}
	if res == nil {
		log.Debug("GetMethods not found", log.Any("reviverTypeName", reviverTypeName))
	}
	return res
}
