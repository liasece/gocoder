package ast

import (
	"go/ast"
	"reflect"
	"regexp"
	"strings"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func (c *CodeDecoder) GetTypeFromASTStructType(ctx DecoderContext, st *ast.StructType) gocoder.Type {
	fields := c.GetStructFieldFromASTStruct(ctx, st)
	res := gocoder.NewStruct(ctx.GetBuildingItemName(), fields)
	res.SetPkg(ctx.GetCurrentPkg())
	return res
}

func (c *CodeDecoder) getTypeFromASTNodeWithName(ctx DecoderContext, st ast.Node) gocoder.Type {
	switch t := st.(type) {
	case *ast.Ident:
		return c.GetTypeFromASTIdent(ctx, t)
	case *ast.StarExpr:
		res := c.getTypeFromASTNodeWithName(ctx, t.X)
		if res != nil {
			return res.TackPtr()
		}
		return nil
	case *ast.SelectorExpr:
		// like time.Time
		pkgName := ctx.GetPkgByAlias(t.X.(*ast.Ident).Name)
		return c.GetType(pkgName + "." + t.Sel.Name)
	case *ast.TypeSpec:
		res := c.getTypeFromASTNodeWithName(ctx, t.Type)
		if res != nil {
			res.AddNotes(c.GetNoteFromCommentGroup(ctx, t.Comment, t.Doc)...)
			if res.Name() != t.Name.Name {
				res = res.WarpNamed(t.Name.Name)
				res.SetPkg(ctx.GetCurrentPkg())
			}
		}
		return res
	case *ast.StructType:
		return c.GetTypeFromASTStructType(ctx, t)
	case *ast.ArrayType:
		res := c.getTypeFromASTNodeWithName(ctx, t.Elt)
		if res != nil {
			return res.Slice()
		}
		return nil
	case *ast.MapType:
		key := c.getTypeFromASTNodeWithName(ctx, t.Key)
		if key == nil {
			return nil
		}
		value := c.getTypeFromASTNodeWithName(ctx, t.Value)
		if value == nil {
			return nil
		}
		return gocoder.NewType(reflect.MapOf(key.RefType(), value.RefType()))
	case *ast.InterfaceType:
		res := c.GetInterfaceFromASTInterfaceType(ctx, t)
		return res
	default:
		log.Warn("name == typeName but type unknown", log.Any("name", ctx.GetBuildingItemName()), log.Any("type", reflect.TypeOf(t)))
	}
	return nil
}

func (c *CodeDecoder) GetTypeFromASTNode(ctx DecoderContext, st ast.Node) gocoder.Type {
	return c.getTypeFromASTNodeWithName(ctx, st)
}

func (c *CodeDecoder) GetTypeFromASTIdent(ctx DecoderContext, st *ast.Ident) gocoder.Type {
	typeStr := st.Name
	res := TypeStringToZeroInterface(typeStr)
	if res == nil {
		// not basic type
		if ctx != nil && ctx.GetCurrentPkg() != "" {
			typeStr = ctx.GetCurrentPkg() + "." + typeStr
		}
		t := c.GetType(typeStr)
		if t != nil {
			res = t
		}
	}
	if res == nil {
		return nil
	}
	return res
}

func (c *CodeDecoder) SearchTypeNames(typePkg string, typeNameRegStr string) []string {
	var res []string
	if typeNameRegStr == "" {
		return res
	}
	log.Debug("SearchTypes", log.Any("typeNameRegStr", typeNameRegStr))

	typeNameReg := regexp.MustCompile(`^\**` + typeNameRegStr + `$`)
	for _, pkgV := range c.pkgs.List {
		if typePkg != "" {
			if typePkg != pkgV.Name && typePkg != pkgV.Alias {
				continue
			}
		}
		for _, astFile := range pkgV.Package.Files {
			for _, astDecl := range astFile.Decls {
				if astGenDecl, ok := astDecl.(*ast.GenDecl); ok {
					ast.Walk(walker(func(node ast.Node) bool {
						if node == nil {
							return true
						}
						ts, ok := node.(*ast.TypeSpec)
						if !ok {
							return true
						}
						if typeNameReg.MatchString(ts.Name.Name) {
							res = append(res, ts.Name.Name)
						}
						return true
					}), astGenDecl)
				}
			}
		}
	}
	return res
}

func (c *CodeDecoder) GetType(fullTypeName string) gocoder.Type {
	if fullTypeName == "" {
		return nil
	}
	log.Debug("GetType", log.Any("fullTypeName", fullTypeName))
	if c.DecodedTypes[fullTypeName] != nil {
		return c.DecodedTypes[fullTypeName]
	}
	astLoadedType := &LoadedType{
		Type: gocoder.NewTypeName(fullTypeName),
	}
	c.DecodedTypes[fullTypeName] = astLoadedType

	var resType gocoder.Type
	basicType := TypeStringToZeroInterface(fullTypeName)
	if basicType != nil {
		resType = basicType
	}
	if resType == nil {
		typeTypeName := fullTypeName
		typePkg := ""
		if index := strings.LastIndex(fullTypeName, "."); index > 0 && index < len(fullTypeName)-1 {
			typePkg = fullTypeName[:index]
			typeTypeName = fullTypeName[index+1:]
		}
		for _, pkgV := range c.pkgs.List {
			if typePkg != "" && typePkg != pkgV.Name && typePkg != pkgV.Alias {
				continue
			}
			for _, astFile := range pkgV.Package.Files {
				for _, astDecl := range astFile.Decls {
					if astGenDecl, ok := astDecl.(*ast.GenDecl); ok {
						ast.Walk(walker(func(node ast.Node) bool {
							if node == nil {
								return true
							}
							ts, ok := node.(*ast.TypeSpec)
							if !ok {
								return true
							}
							if ts.Name.Name == typeTypeName {
								// log.Info("found type", log.Any("type", ts.Name.Name), log.Any("pkg", pkgV.Name), log.Any("file", astFile.Name.Name))
								ctx := NewDecoderContextByAstFile(pkgV.Name, typeTypeName, astFile)
								resType = c.GetTypeFromASTNode(ctx, ts)
								if resType != nil {
									if resType.IsStruct() && resType.Package() == "" {
										resType.SetPkg(pkgV.Name)
									}
									resType.AddNotes(c.GetNoteFromCommentGroup(ctx, astGenDecl.Doc)...)
									return false
								}
							}
							return true
						}), astGenDecl)
					}
				}
				if resType != nil {
					break
				}
			}
			if resType != nil {
				break
			}
		}
	}
	if resType == nil {
		c.DecodedTypes[fullTypeName] = nil
	} else {
		astLoadedType.Type = resType
	}
	return resType
}
