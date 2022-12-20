package gocoder

import (
	"fmt"
	"path"
	"strings"
)

type tPkgTool struct {
	// package path to local alias map for tracking imports
	imports map[string]string
}

func NewDefaultPkgTool() PkgTool {
	return &tPkgTool{
		imports: make(map[string]string),
	}
}

// PkgAlias creates and returns and import alias for a given package.
func (m *tPkgTool) PkgAlias(pkgPath string) string {
	return m.pkgAlias(pkgPath)
}

func (m *tPkgTool) SetPkgAlias(pkgPath string, alias string) {
	m.imports[pkgPath] = alias
}

func (m *tPkgTool) PkgAliasMap() map[string]string {
	return m.imports
}

// fixes vendored paths
func fixPkgPathVendoring(pkgPath string) string {
	const vendor = "/vendor/"
	if i := strings.LastIndex(pkgPath, vendor); i != -1 {
		return pkgPath[i+len(vendor):]
	}
	return pkgPath
}

func fixAliasName(alias string) string {
	alias = strings.ReplaceAll(
		strings.ReplaceAll(alias, ".", "_"),
		"-",
		"_",
	)

	if alias[0] == 'v' { // to void conflicting with var names, say v1
		alias = "_" + alias
	}
	return alias
}

// pkgAlias creates and returns and import alias for a given package.
func (m *tPkgTool) pkgAlias(pkgPath string) string {
	if pkgPath == "" {
		return ""
	}
	pkgPath = fixPkgPathVendoring(pkgPath)
	if alias := m.imports[pkgPath]; alias != "" {
		return alias
	}

	for i := 0; ; i++ {
		alias := fixAliasName(path.Base(pkgPath))
		if i > 0 {
			alias += fmt.Sprint(i)
		}

		exists := false
		for _, v := range m.imports {
			if v == alias {
				exists = true
				break
			}
		}

		if !exists {
			m.imports[pkgPath] = alias
			return alias
		}
	}
}
