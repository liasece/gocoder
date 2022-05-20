package cde

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
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

func TypeStringToZeroInterface(str string) gocoder.Type {
	switch str {
	case "bool":
		return gocoder.MustToType(false)
	case "byte":
		return gocoder.MustToType(byte(0))
	case "int":
		return gocoder.MustToType(int(0))
	case "int8":
		return gocoder.MustToType(int8(0))
	case "int16":
		return gocoder.MustToType(int16(0))
	case "int32":
		return gocoder.MustToType(int32(0))
	case "int64":
		return gocoder.MustToType(int64(0))
	case "uint":
		return gocoder.MustToType(uint(0))
	case "uint8":
		return gocoder.MustToType(uint8(0))
	case "uint16":
		return gocoder.MustToType(uint16(0))
	case "uint32":
		return gocoder.MustToType(uint32(0))
	case "uint64":
		return gocoder.MustToType(uint64(0))
	case "float32":
		return gocoder.MustToType(float32(0))
	case "float64":
		return gocoder.MustToType(float64(0))
	case "complex64":
		return gocoder.MustToType(complex64(0))
	case "complex128":
		return gocoder.MustToType(complex128(0))
	case "string":
		return gocoder.MustToType("")
	case "time.Time":
		return gocoder.MustToType(time.Time{})

	case "[]bool":
		return gocoder.MustToType([]bool{})
	case "[]int":
		return gocoder.MustToType([]int{})
	case "[]int8":
		return gocoder.MustToType([]int8{})
	case "[]int16":
		return gocoder.MustToType([]int16{})
	case "[]int32":
		return gocoder.MustToType([]int32{})
	case "[]int64":
		return gocoder.MustToType([]int64{})
	case "[]uint":
		return gocoder.MustToType([]uint{})
	case "[]uint8":
		return gocoder.MustToType([]uint8{})
	case "[]uint16":
		return gocoder.MustToType([]uint16{})
	case "[]uint32":
		return gocoder.MustToType([]uint32{})
	case "[]uint64":
		return gocoder.MustToType([]uint64{})
	case "[]float32":
		return gocoder.MustToType([]float32{})
	case "[]float64":
		return gocoder.MustToType([]float64{})
	case "[]complex64":
		return gocoder.MustToType([]complex64{})
	case "[]complex128":
		return gocoder.MustToType([]complex128{})
	case "[]string":
		return gocoder.MustToType([]string{})
	case "[]time.Time":
		return gocoder.MustToType([]time.Time{})

	case "interface{}":
		return gocoder.MustToType(InterfaceType)
	}
	return nil
}

type ASTCoder struct {
	fset *token.FileSet
	pkgs map[string]ast.Node
}

func (c *ASTCoder) loadTypeFromASTStructFields(st *ast.StructType) ([]gocoder.Field, error) {
	fields := make([]gocoder.Field, 0)
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
		typ := TypeStringToZeroInterface(typeStr)
		if typ == nil {
			t, err := c.loadTypeFromSourceFileSet(typeStr)
			if err != nil {
				log.Error("loadTypeFromASTStructFields loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr))
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
				}
				typ = t
			}
		}
		if typ == nil {
			log.Warn("not found type", log.Any("typeStr", typeStr), log.Any("st", st))
		} else if !('a' <= name[0] && name[0] <= 'z' || name[0] == '_') {
			if isSlice {
				typ = typ.Slice()
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

func (c *ASTCoder) loadTypeFromASTStructType(st *ast.StructType) (gocoder.Type, error) {
	fields, err := c.loadTypeFromASTStructFields(st)
	if err != nil {
		return nil, err
	}
	return gocoder.NewStruct("", fields).GetType(), nil
}

func (c *ASTCoder) loadTypeFromASTIdent(st *ast.Ident) (gocoder.Type, error) {
	typeStr := st.Name
	// log.Error("walker find", log.Reflect("info", st), log.Any("type", reflect.TypeOf(st)))
	res := TypeStringToZeroInterface(typeStr)
	if res == nil {
		t, err := c.loadTypeFromSourceFileSet(typeStr)
		if err != nil {
			log.Warn("loadTypeFromASTIdent loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("typeStr", typeStr), log.Any("obj", st.Obj), log.Any("st", st))
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

func (c *ASTCoder) loadTypeFromSourceFileSet(typeName string) (gocoder.Type, error) {
	var resType gocoder.Type
	var resErr error
	typeTypeName := typeName
	// typePkgName := ""
	if ss := strings.Split(typeName, "."); len(ss) == 2 {
		// typePkgName = ss[0]
		typeTypeName = ss[1]
	}
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
			if ts.Name.Name != typeTypeName {
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
		// 	return nil, errors.New("not found type: " + typeName + "(" + typeTypeName + ")")
		log.Error("loadTypeFromSourceFileSet not found type", log.Any("typeName", typeName), log.Any("typeTypeName", typeTypeName))
	}
	return resType, resErr
}

func ParseDir(fset *token.FileSet, path string, filter func(fs.FileInfo) bool, mode parser.Mode) (pkgs map[string]*ast.Package, first error) {
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	pkgs = make(map[string]*ast.Package)
	for _, d := range list {
		if d.IsDir() {
			ps, err := ParseDir(fset, filepath.Join(path, d.Name()), filter, mode)
			if err != nil {
				return nil, err
			}
			for k, v := range ps {
				pkgs[k] = v
			}
			continue
		}
		if !strings.HasSuffix(d.Name(), ".go") {
			continue
		}
		if filter != nil {
			info, err := d.Info()
			if err != nil {
				return nil, err
			}
			if !filter(info) {
				continue
			}
		}
		filename := filepath.Join(path, d.Name())
		if src, err := parser.ParseFile(fset, filename, nil, mode); err == nil {
			name := src.Name.Name
			pkg, found := pkgs[name]
			if !found {
				pkg = &ast.Package{
					Name:  name,
					Files: make(map[string]*ast.File),
				}
				pkgs[name] = pkg
			}
			pkg.Files[filename] = src
		} else if first == nil {
			first = err
		}
	}

	return
}

func LoadTypeFromSource(path string, typeName string) (gocoder.Type, error) {
	fset := token.NewFileSet()
	ps := make(map[string]ast.Node)
	pathSS := strings.Split(path, ",")
	for _, path := range pathSS {
		if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
			// from path
			pkgs, err := ParseDir(fset, path, nil, parser.AllErrors)
			if err != nil {
				return nil, err
			}
			for k, v := range pkgs {
				ps[k] = v
			}
		} else {
			node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return nil, err
			}
			ps[node.Name.Name] = node
		}
	}
	c := &ASTCoder{
		fset: fset,
		pkgs: ps,
	}
	return c.loadTypeFromSourceFileSet(typeName)
}
