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
			ast.Print(c.fset, node)
			if ts.Recv == nil || len(ts.Recv.List) == 0 {
				return true
			}
			reviverType, err := c.GetTypeFromASTExpr(ts.Recv.List[0].Type, opt)
			if err != nil {
				log.Error("GetMethods GetTypeFromASTExpr error", log.Any("err", err))
				return true
			}
			if reviverType.String() != reviverTypeName && reviverType.String() != "*"+reviverTypeName {
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
		log.Warn("GetTypeFromSourceFileSet not found type", log.Any("reviverTypeName", reviverTypeName))
	}
	return res, nil
}
