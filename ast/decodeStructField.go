package ast

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func (c *ASTCoder) GetStructFieldFromASTStruct(st *ast.StructType, opt *gocoder.ToCodeOption) ([]gocoder.Field, error) {
	fields := make([]gocoder.Field, 0)
	// log.Error("walker find", log.Any("info", st.Fields), log.Any("type", reflect.TypeOf(st.Fields)))
	afterHandleFuncPtr := func(typ gocoder.Type) gocoder.Type {
		return typ.TackPtr()
	}
	afterHandleFuncSlice := func(typ gocoder.Type) gocoder.Type {
		return typ.Slice()
	}
	for _, arg := range st.Fields.List {
		typeAfterHandle := make([]func(typ gocoder.Type) gocoder.Type, 0)
		name := ""
		if len(arg.Names) > 0 {
			name = arg.Names[0].Name
		}
		typeStr := ""
		astType := arg.Type
		for {
			if se, ok := astType.(*ast.StarExpr); ok {
				// 指针
				astType = se.X
				typeAfterHandle = append(typeAfterHandle, afterHandleFuncPtr)
			} else {
				break
			}
		}

		pkg := ""
		if se, ok := astType.(*ast.SelectorExpr); ok {
			// typeStr = se.X.(*ast.Ident).Name + "." + se.Sel.Name
			typeStr = se.Sel.Name
			pkg = se.X.(*ast.Ident).Name
			// log.Error("walker find", log.Reflect("info", arg.Tag.Value), log.Any("type", reflect.TypeOf(arg.Type)), log.Any("namesLen", len(arg.Names)), log.Reflect("se", se), log.Any("seType", reflect.TypeOf(se.Sel)))
		}
		if se, ok := astType.(*ast.Ident); ok {
			typeStr = se.Name
		}
		if se, ok := astType.(*ast.ArrayType); ok {
			if tmp, ok := se.Elt.(*ast.StarExpr); ok {
				typeStr = tmp.X.(*ast.Ident).Name
				typeAfterHandle = append(typeAfterHandle, afterHandleFuncPtr)
			} else {
				typeStr = se.Elt.(*ast.Ident).Name
			}
			typeAfterHandle = append(typeAfterHandle, afterHandleFuncSlice)
		}
		if _, ok := astType.(*ast.InterfaceType); ok {
			typeStr = "interface{}"
		}
		if typeStr == "" {
			err := ast.Print(c.fset, astType)
			if err != nil {
				log.Error("typeStr is empty, Print error", log.ErrorField(err), log.Any("astType", astType))
			}
		}
		typ := TypeStringToZeroInterface(toTypeStr(pkg, typeStr))
		if typ == nil {
			t, err := c.GetType(toTypeStr(pkg, typeStr), opt)
			if err != nil {
				log.Error("GetTypeFromASTStructFields GetTypeFromSourceFileSet error", log.ErrorField(err), log.Any("pkg", pkg), log.Any("typeStr", typeStr))
				return nil, err
			}
			if t != nil {
				if name == "" {
					// 匿名成员结构体
					for i := 0; i < t.NumField(); i++ {
						fields = append(fields, t.Field(i))
					}
					continue
				} else if t.Kind() == reflect.Struct {
					// log.Info("skip type: Struct", log.Any("typeStr", typeStr))
					// continue
					t.SetNamed(typeStr)
					if t.Package() == "" {
						t.SetPkg(pkg)
						log.Warn("GetTypeFromASTStructFields set pkg", log.Any("name", t.Name()), log.Any("str", t.String()), log.Any("pkg", t.Package()))
					}
					// log.Debug("in SetNamed", log.Any("name", t.Name()), log.Any("str", t.String()), log.Any("pkg", t.Package()))
				}
				typ = t
			}
		}
		if typ == nil {
			log.Warn("not found type", log.Any("typeStr", typeStr), log.Any("st", st))
		} else if !('a' <= name[0] && name[0] <= 'z' || name[0] == '_') {
			for _, f := range typeAfterHandle {
				typ = f(typ)
			}
			var tag string
			if arg.Tag != nil {
				tag = strings.ReplaceAll(arg.Tag.Value, "`", "")
			}
			fields = append(fields, gocoder.NewField(name, typ, tag))
		} else {
			log.Info("skip type", log.Any("typeStr", typeStr), log.Any("st", st))
		}
	}
	return fields, nil
}
