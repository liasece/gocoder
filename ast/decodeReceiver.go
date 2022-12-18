package ast

import (
	"go/ast"

	"github.com/liasece/gocoder"
)

func (c *ASTCoder) GetReceiverFromASTField(st *ast.Field, opt *gocoder.ToCodeOption) (gocoder.Receiver, error) {
	var name string
	if len(st.Names) > 0 {
		name = st.Names[0].Name
	}
	t, err := c.GetTypeFromASTExpr(st.Type, opt)
	if err != nil {
		return nil, err
	}
	return gocoder.NewReceiver(name, t), nil
}
