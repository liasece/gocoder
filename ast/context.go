package ast

import (
	"go/ast"
	"strings"
)

type DecoderContext interface {
	GetPkgByAlias(alias string) string
	GetBuildingItemName() string
	GetCurrentPkg() string
}

type decoderContext struct {
	currentPkg       string
	buildingItemName string
	pkgAliasMap      map[string]string // key: imported package alias, value: imported package full path
}

func NewDecoderContextByAstFile(currentPkg string, buildingItemName string, astFile *ast.File) *decoderContext {
	pkgAliasMap := make(map[string]string) // key: imported package alias, value: imported package full path
	for _, imp := range astFile.Imports {
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		pkg := strings.ReplaceAll(imp.Path.Value, "\"", "")
		if alias == "" || alias == "_" {
			ss := strings.Split(pkg, "/")
			if len(ss) > 0 {
				alias = ss[len(ss)-1]
			}
		}
		pkgAliasMap[alias] = pkg
	}
	res := &decoderContext{
		currentPkg:       currentPkg,
		buildingItemName: buildingItemName,
		pkgAliasMap:      pkgAliasMap,
	}
	return res
}

func (c *decoderContext) GetPkgByAlias(alias string) string {
	if c.pkgAliasMap == nil {
		return alias
	}
	path, ok := c.pkgAliasMap[alias]
	if !ok {
		return alias
	}
	return path
}

func (c *decoderContext) GetBuildingItemName() string {
	return c.buildingItemName
}

func (c *decoderContext) GetCurrentPkg() string {
	return c.currentPkg
}
