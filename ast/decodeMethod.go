package ast

import (
	"go/ast"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func (c *ASTCoder) GetMethods(reviverTypeName string, opt *gocoder.ToCodeOption) ([]gocoder.Func, error) {
	var res []gocoder.Func
	for _, pkgV := range c.pkgs {
		_, node := pkgV.name, pkgV.node
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
			reviverType, err := c.GetTypeFromASTNode(ts.Recv.List[0].Type, opt)
			if err != nil {
				log.Error("GetMethods GetTypeFromASTNode error", log.Any("err", err))
				return true
			}
			name := ""
			if reviverType != nil {
				name = reviverType.String()
			}
			if name != reviverTypeName && name != "*"+reviverTypeName {
				return true
			}
			fn, err := c.GetFuncsFromASTFuncDecl(ts, opt)
			if err != nil {
				log.Error("GetMethods GetFuncsFromASTField error", log.Any("err", err))
				return true
			}
			res = append(res, fn)
			return true
		}), node)
	}
	if res == nil {
		log.Debug("GetMethods not found", log.Any("reviverTypeName", reviverTypeName))
	}
	return res, nil
}
