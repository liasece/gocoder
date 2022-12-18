package ast

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func (c *ASTCoder) GetTypeFromASTStructType(name string, st *ast.StructType, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	fields, err := c.GetStructFieldFromASTStruct(st, opt)
	if err != nil {
		return nil, err
	}
	res := gocoder.NewStruct(name, fields).GetType()
	for i := 0; i < res.NumField(); i++ {
		f := res.Field(i)
		t := f.GetType()
		if t.Package() != res.Package() && c.importPkgs[t.Package()] != "" {
			t.SetPkg(c.importPkgs[t.Package()])
		}
	}
	return res, nil
}

func (c *ASTCoder) GetTypeFromASTExpr(st ast.Expr, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	switch t := st.(type) {
	case *ast.Ident:
		return c.GetTypeFromASTIdent(st.(*ast.Ident), opt)
	case *ast.StarExpr:
		res, err := c.GetTypeFromASTExpr(st.(*ast.StarExpr).X, opt)
		if err != nil {
			return nil, err
		}
		return res.TackPtr(), nil
	case *ast.SelectorExpr:
		str := t.X.(*ast.Ident).Name + "." + t.Sel.Name
		return TypeStringToZeroInterface(str), nil
	}
	return nil, nil
}

func (c *ASTCoder) GetTypeFromASTIdent(st *ast.Ident, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	typeStr := st.Name
	// log.Warn("GetTypeFromASTIdent walker find", log.Any("typeStr", typeStr), log.Reflect("info", st), log.Any("type", reflect.TypeOf(st)))
	res := TypeStringToZeroInterface(typeStr)
	if res == nil {
		// not basic type
		t, err := c.GetType(typeStr, opt)
		if err != nil {
			log.Warn("GetTypeFromASTIdent GetTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr), log.Any("obj", st.Obj), log.Any("st", st))
			ast.Print(c.fset, st)
			return nil, nil
		}
		if t == nil {
			log.Warn("not found type", log.ErrorField(err), log.Any("typeStr", typeStr), log.Any("obj", st.Obj), log.Any("st", st))
		} else {
			res = t
		}
	}
	if res == nil {
		return nil, nil
	}
	return res, nil
}

func (c *ASTCoder) GetType(typeName string, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	var resType gocoder.Type
	var resErr error
	typeTypeName := typeName
	// typePkgName := ""
	if ss := strings.Split(typeName, "."); len(ss) == 2 {
		// typePkgName = ss[0]
		typeTypeName = ss[1]
	}
	for _, pkgV := range c.pkgs {
		pkg, node := pkgV.name, pkgV.node
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
					// log.Error("add import", log.Any("key", ss[len(ss)-1]), log.Any("value", path))
					// ast.Print(c.fset, ts)
				}
			}
			// log.Error("walker", log.Any("node", node), log.Any("nodeType", reflect.TypeOf(node)))
			ts, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if ts.Name.Name == typeTypeName {
				// found target type
				if st, ok := ts.Type.(*ast.StructType); ok {
					resType, resErr = c.GetTypeFromASTStructType(ts.Name.Name, st, opt)
					if pkgPath := opt.GetPkgPath(); pkgPath != nil && pkgInReference(*pkgPath) == pkg {
						resType.SetPkg(*pkgPath)
						// log.Warn("GetTypeFromSourceFileSet GetTypeFromASTStructType set pkg from opt", log.Reflect("ts.Name.Name", ts.Name.Name), log.Any("pkgPath", pkgPath), log.Any("nowPkg", pkg))
					} else {
						resType.SetPkg(pkg)
					}
					return false
				}
				if st, ok := ts.Type.(*ast.Ident); ok {
					resType, resErr = c.GetTypeFromASTIdent(st, opt)
					return false
				}
				log.Error("ts.Name.Name == typeName but type unknown", log.Any("loadTypeName", typeName), log.Any("typeTypeName", typeTypeName), log.Any("type", reflect.TypeOf(ts.Type)))
			}
			return true
		}), node)
		if resErr != nil || resType != nil {
			break
		}
	}
	if resType == nil {
		// 	return nil, errors.New("not found type: " + typeName + "(" + typeTypeName + ")")
		log.Warn("GetTypeFromSourceFileSet not found type", log.Any("typeName", typeName), log.Any("typeTypeName", typeTypeName))
	}
	return resType, resErr
}
