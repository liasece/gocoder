package ast

import (
	"go/ast"

	"github.com/liasece/gocoder"
)

func (c *CodeDecoder) GetReceiverFromASTField(ctx DecoderContext, st *ast.Field) gocoder.Receiver {
	var name string
	if len(st.Names) > 0 {
		name = st.Names[0].Name
	}
	t := c.getTypeFromASTNodeWithName(ctx, st.Type)
	res := gocoder.NewReceiver(name, t)
	res.AddNotes(c.GetNoteFromCommentGroup(ctx, st.Doc, st.Comment)...)
	return res
}
