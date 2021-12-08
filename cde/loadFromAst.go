package cde

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"time"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

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

func TypeStringToZeroInterface(str string) interface{} {
	switch str {
	case "bool":
		return false
	case "int":
		return int(0)
	case "int8":
		return int8(0)
	case "int16":
		return int16(0)
	case "int32":
		return int32(0)
	case "int64":
		return int64(0)
	case "uint":
		return uint(0)
	case "uint8":
		return uint8(0)
	case "uint16":
		return uint16(0)
	case "uint32":
		return uint32(0)
	case "uint64":
		return uint64(0)
	case "float32":
		return float32(0)
	case "float64":
		return float64(0)
	case "complex64":
		return complex64(0)
	case "complex128":
		return complex128(0)
	case "string":
		return ""
	case "time.Time":
		return time.Time{}
	}
	return nil
}

type ASTCoder struct {
	fset *token.FileSet
	node ast.Node
}

func (c *ASTCoder) loadTypeFromASTStructType(st *ast.StructType) (gocoder.Type, error) {
	fields := make([]reflect.StructField, 0)
	// log.Error("walker find", log.Any("info", st.Fields), log.Any("type", reflect.TypeOf(st.Fields)))
	for _, arg := range st.Fields.List {
		name := ""
		if len(arg.Names) > 0 {
			name = arg.Names[0].Name
		}
		typeStr := ""
		if se, ok := arg.Type.(*ast.SelectorExpr); ok {
			typeStr = se.X.(*ast.Ident).Name + "." + se.Sel.Name
			// log.Error("walker find", log.Reflect("info", arg.Tag.Value), log.Any("type", reflect.TypeOf(arg.Type)), log.Any("namesLen", len(arg.Names)), log.Reflect("se", se), log.Any("seType", reflect.TypeOf(se.Sel)))
		}
		if se, ok := arg.Type.(*ast.Ident); ok {
			typeStr = se.Name
		}
		zeroI := TypeStringToZeroInterface(typeStr)
		if zeroI == nil {
			t, err := c.loadTypeFromSourceFileSet(typeStr)
			if err != nil {
				log.Error("loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr))
				return nil, err
			}
			zeroI = reflect.Zero(t.RefType()).Interface()
		}
		fields = append(fields, reflect.StructField{
			Name: name,
			Type: reflect.TypeOf(zeroI),
			Tag:  reflect.StructTag(strings.ReplaceAll(arg.Tag.Value, "`", "")),
		})
	}
	return Type(reflect.StructOf(fields)), nil
}

func (c *ASTCoder) loadTypeFromASTIdent(st *ast.Ident) (gocoder.Type, error) {
	typeStr := st.Name
	// log.Error("walker find", log.Reflect("info", st), log.Any("type", reflect.TypeOf(st)))
	zeroI := TypeStringToZeroInterface(typeStr)
	if zeroI == nil {
		t, err := c.loadTypeFromSourceFileSet(typeStr)
		if err != nil {
			log.Error("loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr))
			return nil, err
		}
		zeroI = reflect.Zero(t.RefType()).Interface()
	}
	return Type(zeroI), nil
}

func (c *ASTCoder) loadTypeFromSourceFileSet(typeName string) (gocoder.Type, error) {
	var resType gocoder.Type
	var resErr error
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
	}), c.node)
	if resType == nil {
		return nil, errors.New("not found type: " + typeName)
	}
	return resType, resErr
}

func LoadTypeFromSource(path string, typeName string) (gocoder.Type, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	c := &ASTCoder{
		fset: fset,
		node: node,
	}
	return c.loadTypeFromSourceFileSet(typeName)
}
