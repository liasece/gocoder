package gocoder

import (
	"fmt"
	"reflect"
	"strings"
)

func useOfCastType(t Type, tool PkgTool) string {
	return typeStringOut(t, tool, "")
}

func typeStrIter(refType reflect.Type, path string, tool PkgTool) string {
	if refType == nil {
		return ""
	}
	if refType.PkgPath() != "" {
		return path + tool.PkgAlias(refType.PkgPath()) + "." + refType.Name()
	}
	switch refType.Kind() {
	case reflect.Ptr:
		return typeStrIter(refType.Elem(), path+"*", tool)
	case reflect.Slice:
		return typeStrIter(refType.Elem(), path+"[]", tool)
	case reflect.Map:
		return typeStrIter(refType.Elem(), path+"map["+typeStrIter(refType.Key(), "", tool)+"]", tool)
	default:
		return refType.String()
	}
}

func typeStringOut(t Type, tool PkgTool, toPkg string) string {
	str := t.CurrentCode()
	if tool != nil {
		if pkg := t.Package(); pkg != "" && pkg != toPkg {
			if str == "" {
				return ""
			}
			prefix := ""
			if strings.HasPrefix(str, "[]") {
				str = str[2:]
				prefix = "[]"
			}
			if strings.HasPrefix(str, "*") {
				str = str[1:]
				prefix = "*"
			}
			return prefix + tool.PkgAlias(pkg) + "." + str
		}
		if sl := strings.Split(str, "."); len(sl) == 2 {
			path := ""
			{
				path = typeStrIter(t.RefType(), "", tool)
			}
			if path != "" {
				return path
			}
		}
	}
	return str
}

func getZeroValueCode(t Type, tool PkgTool) string {
	originT := t
	getPtr := ""
	if t == nil || t.RefType() == nil {
		return `nil`
	}
	if t.Kind() == reflect.Ptr {
		getPtr = "&"
		t = t.Elem()
		if t.Kind() != reflect.Struct {
			return `nil`
		}
	}
	switch t.Kind() {
	case reflect.Struct:
		return getPtr + useOfCastType(t, tool) + "{}"
	case reflect.Slice, reflect.Map:
		if getPtr != "" {
			return "(*" + useOfCastType(t, tool) + ")" + "(nil)"
		}
		return "(" + useOfCastType(t, tool) + ")" + "(nil)"
	default:
		zero := fmt.Sprint(reflect.Zero(t.RefType()).Interface())
		if t.Kind() == reflect.String {
			zero = `""`
		}
		if getPtr != "" {
			return "func() " + useOfCastType(originT, tool) + " { v := (" + useOfCastType(t, tool) + ")(" + zero + "); return &v }()"
		}
		return zero
	}
}
