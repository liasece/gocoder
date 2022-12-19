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

func (c *ASTCoder) getTypeFromASTNodeWithName(name string, st ast.Node, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	switch t := st.(type) {
	case *ast.Ident:
		return c.GetTypeFromASTIdent(t, opt)
	case *ast.StarExpr:
		res, err := c.getTypeFromASTNodeWithName(name, t.X, opt)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res.TackPtr(), nil
		} else {
			return nil, nil
		}
	case *ast.SelectorExpr:
		str := t.X.(*ast.Ident).Name + "." + t.Sel.Name
		return TypeStringToZeroInterface(str), nil
	case *ast.TypeSpec:
		return c.getTypeFromASTNodeWithName(t.Name.Name, t.Type, opt)
	case *ast.StructType:
		return c.GetTypeFromASTStructType(name, t, opt)
	case *ast.ArrayType:
		res, err := c.getTypeFromASTNodeWithName(name, t.Elt, opt)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res.Slice(), nil
		} else {
			return nil, nil
		}
	case *ast.MapType:
		key, err := c.getTypeFromASTNodeWithName(name, t.Key, opt)
		if err != nil {
			return nil, err
		}
		if key == nil {
			return nil, nil
		}
		value, err := c.getTypeFromASTNodeWithName(name, t.Value, opt)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		return gocoder.NewType(reflect.MapOf(key.RefType(), value.RefType())), nil
	case *ast.InterfaceType:
		res, err := c.GetInterfaceFromASTInterfaceType(name, t, opt)
		if err != nil {
			return nil, err
		}
		return res.GetType(), nil
	default:
		log.Warn("name == typeName but type unknown", log.Any("name", name), log.Any("type", reflect.TypeOf(t)))
	}
	return nil, nil
}

func (c *ASTCoder) GetTypeFromASTNode(st ast.Node, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	return c.getTypeFromASTNodeWithName("", st, opt)
}

func (c *ASTCoder) GetTypeFromASTIdent(st *ast.Ident, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	typeStr := st.Name
	res := TypeStringToZeroInterface(typeStr)
	if res == nil {
		// not basic type
		t, err := c.GetType(typeStr, opt)
		if err != nil {
			log.Warn("GetTypeFromASTIdent GetTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr), log.Any("obj", st.Obj), log.Any("st", st))
			return nil, nil
		}
		if t != nil {
			res = t
		}
	}
	if res == nil {
		return nil, nil
	}
	return res, nil
}

func (c *ASTCoder) GetType(typeName string, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	if typeName == "" {
		return nil, nil
	}
	if c.DecodedTypes[typeName] != nil {
		return c.DecodedTypes[typeName], nil
	}
	astTyped := &ASTTyped{
		Type: gocoder.NewTypeName(typeName),
	}
	c.DecodedTypes[typeName] = astTyped

	var resType gocoder.Type
	var resErr error
	typeTypeName := typeName
	if ss := strings.Split(typeName, "."); len(ss) == 2 {
		typeTypeName = ss[1]
	}
	basicType := TypeStringToZeroInterface(typeTypeName)
	if basicType != nil {
		resType = basicType
	} else {
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
					}
				}
				ts, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}
				if ts.Name.Name == typeTypeName {
					resType, resErr = c.GetTypeFromASTNode(ts, opt)
					if resType != nil && resType.IsStruct() {
						if pkgPath := opt.GetPkgPath(); pkgPath != nil && pkgInReference(*pkgPath) == pkg {
							resType.SetPkg(*pkgPath)
						} else {
							resType.SetPkg(pkg)
						}
					}
				}
				return true
			}), node)
			if resErr != nil || resType != nil {
				break
			}
		}
	}
	if resType == nil {
		c.DecodedTypes[typeName] = nil
	} else {
		astTyped.Type = resType
	}
	return resType, resErr
}
