package ast

import (
	"go/ast"
	"strings"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func (c *CodeDecoder) GetStructFieldFromASTStruct(ctx DecoderContext, st *ast.StructType) []gocoder.Field {
	fields := make([]gocoder.Field, 0)
	for _, astField := range st.Fields.List {
		name := ""
		if len(astField.Names) > 0 {
			name = astField.Names[0].Name
		}
		astType := astField.Type

		typ := c.getTypeFromASTNodeWithName(ctx, astType)
		if typ == nil {
			log.Debug("GetStructFieldFromASTStruct not found type", log.Any("astType", astType), log.Any("st", st))
			continue
		}

		if name == "" {
			// 匿名成员结构体
			for i := 0; i < typ.NumField(); i++ {
				fields = append(fields, typ.Field(i))
			}
			continue
		}

		if !('a' <= name[0] && name[0] <= 'z' || name[0] == '_') {
			var tag string
			if astField.Tag != nil {
				tag = strings.ReplaceAll(astField.Tag.Value, "`", "")
			}
			fields = append(fields, gocoder.NewField(name, typ, tag))
		} else {
			log.Info("skip type", log.Any("name", name), log.Any("astType", astType), log.Any("st", st))
		}
	}
	return fields
}
