package ast

import (
	"go/ast"

	"github.com/liasece/gocoder"
)

func (c *CodeDecoder) GetFuncsFromASTFieldList(ctx DecoderContext, receiver gocoder.Receiver, st *ast.FieldList) []gocoder.Func {
	fs := make([]gocoder.Func, 0)
	for _, arg := range st.List {
		f := c.GetFuncFromASTField(ctx, receiver, arg)
		if f != nil {
			fs = append(fs, f)
		}
	}
	return fs
}

func (c *CodeDecoder) GetFuncsFromASTFuncDecl(ctx DecoderContext, st *ast.FuncDecl) gocoder.Func {
	var name string
	{
		// get name
		name = st.Name.Name
	}
	receiver := c.GetReceiverFromASTField(ctx, st.Recv.List[0])
	return c.GetFuncsFromASTFuncType(ctx, receiver, name, st.Type)
}

func (c *CodeDecoder) GetFuncFromASTField(ctx DecoderContext, receiver gocoder.Receiver, st *ast.Field) gocoder.Func {
	var name string
	{
		// get name
		if len(st.Names) > 0 {
			name = st.Names[0].Name
		}
	}
	return c.GetFuncsFromASTFuncType(ctx, receiver, name, st.Type.(*ast.FuncType))
}

func (c *CodeDecoder) GetFuncsFromASTFuncType(ctx DecoderContext, receiver gocoder.Receiver, name string, se *ast.FuncType) gocoder.Func {
	var args []gocoder.Arg
	var returns []gocoder.Type

	if se.Params != nil {
		for _, arg := range se.Params.List {
			argType := c.GetTypeFromASTNode(ctx, arg.Type)
			for _, argName := range arg.Names {
				args = append(args, gocoder.NewArg(argName.Name, argType, false))
			}
		}
	}
	if se.Results != nil {
		for _, arg := range se.Results.List {
			argType := c.GetTypeFromASTNode(ctx, arg.Type)
			returns = append(returns, argType)
		}
	}

	return gocoder.NewFunc(gocoder.FuncTypeDefault, name, receiver, args, returns)
}
