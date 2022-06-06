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
	fset       *token.FileSet
	pkgs       map[string]ast.Node
	importPkgs map[string]string
}

func toTypeStr(pkg string, name string) string {
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

func (c *ASTCoder) loadTypeFromASTStructFields(st *ast.StructType, opt *gocoder.ToCodeOption) ([]gocoder.Field, error) {
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
			t, err := c.loadTypeFromSourceFileSet(toTypeStr(pkg, typeStr), opt)
			if err != nil {
				log.Error("loadTypeFromASTStructFields loadTypeFromSourceFileSet error", log.ErrorField(err), log.Any("pkg", pkg), log.Any("typeStr", typeStr))
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
						log.Warn("loadTypeFromASTStructFields set pkg", log.Any("name", t.Name()), log.Any("str", t.String()), log.Any("pkg", t.Package()))
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

func (c *ASTCoder) loadTypeFromASTStructType(name string, st *ast.StructType, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	fields, err := c.loadTypeFromASTStructFields(st, opt)
	if err != nil {
		return nil, err
	}
	res := gocoder.NewStruct(name, fields).GetType()
	for i := 0; i < res.NumField(); i++ {
		f := res.Field(i)
		t := f.GetType()
		// log.Error("reload loadTypeFromASTStructType", log.Any("pkg", t.Package()), log.Any("resPkg", res.Package()), log.Any("map", c.importPkgs[t.Package()]), log.Any("res", t.String()))
		if t.Package() != res.Package() && c.importPkgs[t.Package()] != "" {
			t.SetPkg(c.importPkgs[t.Package()])
			// log.Warn("set pkg", log.Any("str", t.String()))
			// } else {
			// log.Debug("not set pkg", log.Any("str", t.String()))
		}
	}
	return res, nil
}

func (c *ASTCoder) loadTypeFromASTIdent(st *ast.Ident, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	typeStr := st.Name
	// log.Warn("loadTypeFromASTIdent walker find", log.Any("typeStr", typeStr), log.Reflect("info", st), log.Any("type", reflect.TypeOf(st)))
	res := TypeStringToZeroInterface(typeStr)
	if res == nil {
		// not basic type
		t, err := c.loadTypeFromSourceFileSet(typeStr, opt)
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

func pkgInReference(str string) string {
	ss := strings.Split(str, "/")
	return ss[len(ss)-1]
}

func (c *ASTCoder) loadTypeFromSourceFileSet(typeName string, opt *gocoder.ToCodeOption) (gocoder.Type, error) {
	var resType gocoder.Type
	var resErr error
	typeTypeName := typeName
	// typePkgName := ""
	if ss := strings.Split(typeName, "."); len(ss) == 2 {
		// typePkgName = ss[0]
		typeTypeName = ss[1]
	}
	for pkg, node := range c.pkgs {
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
					resType, resErr = c.loadTypeFromASTStructType(ts.Name.Name, st, opt)
					if pkgPath := opt.GetPkgPath(); pkgPath != nil && pkgInReference(*pkgPath) == pkg {
						resType.SetPkg(*pkgPath)
						// log.Warn("loadTypeFromSourceFileSet loadTypeFromASTStructType set pkg from opt", log.Reflect("ts.Name.Name", ts.Name.Name), log.Any("pkgPath", pkgPath), log.Any("nowPkg", pkg))
					} else {
						resType.SetPkg(pkg)
					}
					return false
				}
				if st, ok := ts.Type.(*ast.Ident); ok {
					resType, resErr = c.loadTypeFromASTIdent(st, opt)
					return false
				}
				log.Error("ts.Name.Name == typeName but type unknown", log.Any("type", reflect.TypeOf(ts.Type)))
			}
			return true
		}), node)
		if resErr != nil || resType != nil {
			break
		}
	}
	if resType == nil {
		// 	return nil, errors.New("not found type: " + typeName + "(" + typeTypeName + ")")
		log.Warn("loadTypeFromSourceFileSet not found type", log.Any("typeName", typeName), log.Any("typeTypeName", typeTypeName))
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

func LoadTypeFromSource(path string, typeName string, opts ...*gocoder.ToCodeOption) (gocoder.Type, error) {
	opt := gocoder.MergeToCodeOpt(opts...)

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
		fset:       fset,
		pkgs:       ps,
		importPkgs: make(map[string]string),
	}
	return c.loadTypeFromSourceFileSet(typeName, opt)
}
