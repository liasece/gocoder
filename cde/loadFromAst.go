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

	case "[]bool":
		return []bool{}
	case "[]int":
		return []int{}
	case "[]int8":
		return []int8{}
	case "[]int16":
		return []int16{}
	case "[]int32":
		return []int32{}
	case "[]int64":
		return []int64{}
	case "[]uint":
		return []uint{}
	case "[]uint8":
		return []uint8{}
	case "[]uint16":
		return []uint16{}
	case "[]uint32":
		return []uint32{}
	case "[]uint64":
		return []uint64{}
	case "[]float32":
		return []float32{}
	case "[]float64":
		return []float64{}
	case "[]complex64":
		return []complex64{}
	case "[]complex128":
		return []complex128{}
	case "[]string":
		return []string{}
	case "[]time.Time":
		return []time.Time{}
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
		if typeStr == "" {
			err := ast.Print(c.fset, astType)
			if err != nil {
				log.Error("typeStr is empty, Print error", log.ErrorField(err), log.Any("astType", astType))
			}
		}
		zeroI := TypeStringToZeroInterface(typeStr)
		if zeroI == nil {
			t, err := c.loadTypeFromSourceFileSet(typeStr)
			if err != nil {
				log.Error("loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr))
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
			zeroI = reflect.Zero(t.RefType()).Interface()
		}
		refType := reflect.TypeOf(zeroI)
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
	zeroI := TypeStringToZeroInterface(typeStr)
	if zeroI == nil {
		t, err := c.loadTypeFromSourceFileSet(typeStr)
		if err != nil {
			log.Error("loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr), log.Any("obj", st.Obj))
			return nil, err
		}
		zeroI = reflect.Zero(t.RefType()).Interface()
	}
	return Type(zeroI), nil
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
