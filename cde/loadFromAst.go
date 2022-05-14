package cde

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

var InterfaceType reflect.Type

func init() {
	type T struct {
		A interface{}
	}
	InterfaceType = reflect.ValueOf(T{}).Field(0).Type()
}

// wrap a function to fulfill ast.Visitor interface
type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

func TypeStringToReflectKind(str string) reflect.Kind {
	switch str {
	case "string":
		return reflect.String
	}
	return reflect.Invalid
}

func TypeStringToZeroInterface(str string) reflect.Type {
	switch str {
	case "bool":
		return reflect.TypeOf(false)
	case "int":
		return reflect.TypeOf(int(0))
	case "int8":
		return reflect.TypeOf(int8(0))
	case "int16":
		return reflect.TypeOf(int16(0))
	case "int32":
		return reflect.TypeOf(int32(0))
	case "int64":
		return reflect.TypeOf(int64(0))
	case "uint":
		return reflect.TypeOf(uint(0))
	case "uint8":
		return reflect.TypeOf(uint8(0))
	case "uint16":
		return reflect.TypeOf(uint16(0))
	case "uint32":
		return reflect.TypeOf(uint32(0))
	case "uint64":
		return reflect.TypeOf(uint64(0))
	case "float32":
		return reflect.TypeOf(float32(0))
	case "float64":
		return reflect.TypeOf(float64(0))
	case "complex64":
		return reflect.TypeOf(complex64(0))
	case "complex128":
		return reflect.TypeOf(complex128(0))
	case "string":
		return reflect.TypeOf("")
	case "time.Time":
		return reflect.TypeOf(time.Time{})

	case "[]bool":
		return reflect.TypeOf([]bool{})
	case "[]int":
		return reflect.TypeOf([]int{})
	case "[]int8":
		return reflect.TypeOf([]int8{})
	case "[]int16":
		return reflect.TypeOf([]int16{})
	case "[]int32":
		return reflect.TypeOf([]int32{})
	case "[]int64":
		return reflect.TypeOf([]int64{})
	case "[]uint":
		return reflect.TypeOf([]uint{})
	case "[]uint8":
		return reflect.TypeOf([]uint8{})
	case "[]uint16":
		return reflect.TypeOf([]uint16{})
	case "[]uint32":
		return reflect.TypeOf([]uint32{})
	case "[]uint64":
		return reflect.TypeOf([]uint64{})
	case "[]float32":
		return reflect.TypeOf([]float32{})
	case "[]float64":
		return reflect.TypeOf([]float64{})
	case "[]complex64":
		return reflect.TypeOf([]complex64{})
	case "[]complex128":
		return reflect.TypeOf([]complex128{})
	case "[]string":
		return reflect.TypeOf([]string{})
	case "[]time.Time":
		return reflect.TypeOf([]time.Time{})

	case "interface{}":
		return InterfaceType
	}
	return nil
}

type ASTCoder struct {
	fset *token.FileSet
	pkgs map[string]ast.Node
}

func (c *ASTCoder) loadTypeFromASTStructFields(st *ast.StructType) ([]reflect.StructField, error) {
	fields := make([]reflect.StructField, 0)
	// log.Error("walker find", log.Any("info", st.Fields), log.Any("type", reflect.TypeOf(st.Fields)))
	for _, arg := range st.Fields.List {
		name := ""
		isSlice := false
		if len(arg.Names) > 0 {
			name = arg.Names[0].Name
		}
		typeStr := ""
		astType := arg.Type
		for {
			if se, ok := astType.(*ast.StarExpr); ok {
				// 指针
				astType = se.X
			} else {
				break
			}
		}

		if se, ok := astType.(*ast.SelectorExpr); ok {
			typeStr = se.X.(*ast.Ident).Name + "." + se.Sel.Name
			// log.Error("walker find", log.Reflect("info", arg.Tag.Value), log.Any("type", reflect.TypeOf(arg.Type)), log.Any("namesLen", len(arg.Names)), log.Reflect("se", se), log.Any("seType", reflect.TypeOf(se.Sel)))
		}
		if se, ok := astType.(*ast.Ident); ok {
			typeStr = se.Name
		}
		if se, ok := astType.(*ast.ArrayType); ok {
			if tmp, ok := se.Elt.(*ast.StarExpr); ok {
				typeStr = tmp.X.(*ast.Ident).Name
			} else {
				typeStr = se.Elt.(*ast.Ident).Name
			}
			isSlice = true
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
		refType := TypeStringToZeroInterface(typeStr)
		if refType == nil {
			t, err := c.loadTypeFromSourceFileSet(typeStr)
			if err != nil {
				log.Error("loadTypeFromASTStructFields loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr))
				return nil, err
			}
			if name == "" {
				// 匿名成员结构体
				for i := 0; i < t.RefType().NumField(); i++ {
					fields = append(fields, t.RefType().Field(i))
				}
				continue
			} else if t.Kind() == reflect.Struct {
				continue
			}
			refType = t.RefType()
		}
		if isSlice {
			refType = reflect.SliceOf(refType)
		}
		fields = append(fields, reflect.StructField{
			Name: name,
			Type: refType,
			Tag:  reflect.StructTag(strings.ReplaceAll(arg.Tag.Value, "`", "")),
		})
	}
	return fields, nil
}

func (c *ASTCoder) loadTypeFromASTStructType(st *ast.StructType) (gocoder.Type, error) {
	fields, err := c.loadTypeFromASTStructFields(st)
	if err != nil {
		return nil, err
	}
	return Type(reflect.StructOf(fields)), nil
}

func (c *ASTCoder) loadTypeFromASTIdent(st *ast.Ident) (gocoder.Type, error) {
	typeStr := st.Name
	// log.Error("walker find", log.Reflect("info", st), log.Any("type", reflect.TypeOf(st)))
	res := TypeStringToZeroInterface(typeStr)
	if res == nil {
		t, err := c.loadTypeFromSourceFileSet(typeStr)
		if err != nil {
			log.Error("loadTypeFromASTIdent loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr), log.Any("obj", st.Obj), log.Any("st", st))
			ast.Print(c.fset, st)
			return nil, err
		}
		res = t.RefType()
	}
	return Type(res), nil
}

func (c *ASTCoder) loadTypeFromSourceFileSet(typeName string) (gocoder.Type, error) {
	var resType gocoder.Type
	var resErr error
	for _, node := range c.pkgs {
		ast.Walk((walker)(func(node ast.Node) bool {
			if node == nil {
				return true
			}
			// log.Error("walker", log.Any("node", node), log.Any("nodeType", reflect.TypeOf(node)))
			ts, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if ts.Name.Name != typeName {
				return true
			}
			if st, ok := ts.Type.(*ast.StructType); ok {
				resType, resErr = c.loadTypeFromASTStructType(st)
				return false
			}
			if st, ok := ts.Type.(*ast.Ident); ok {
				resType, resErr = c.loadTypeFromASTIdent(st)
				return false
			}
			log.Error("ts.Name.Name == typeName but type unknown", log.Any("type", reflect.TypeOf(ts.Type)))
			return true
		}), node)
		if resErr != nil || resType != nil {
			break
		}
	}
	if resType == nil {
		return nil, errors.New("not found type: " + typeName)
	}
	return resType, resErr
}

func LoadTypeFromSource(path string, typeName string) (gocoder.Type, error) {
	fset := token.NewFileSet()
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
		// from path
		pkgs, err := parser.ParseDir(fset, path, nil, parser.AllErrors)
		if err != nil {
			return nil, err
		}
		ps := make(map[string]ast.Node)
		for k, v := range pkgs {
			ps[k] = v
		}
		c := &ASTCoder{
			fset: fset,
			pkgs: ps,
		}
		return c.loadTypeFromSourceFileSet(typeName)
	}
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	c := &ASTCoder{
		fset: fset,
		pkgs: map[string]ast.Node{node.Name.Name: node},
	}
	return c.loadTypeFromSourceFileSet(typeName)
}
