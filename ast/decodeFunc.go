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
	res := c.GetFuncsFromASTFuncType(ctx, receiver, name, st.Type)
	res.AddNotes(c.GetNoteFromCommentGroup(ctx, st.Doc)...)
	return res
}

func (c *CodeDecoder) GetFuncFromASTField(ctx DecoderContext, receiver gocoder.Receiver, st *ast.Field) gocoder.Func {
	var name string
	{
		// get name
		if len(st.Names) > 0 {
			name = st.Names[0].Name
		}
	}
	res := c.GetFuncsFromASTFuncType(ctx, receiver, name, st.Type.(*ast.FuncType))
	res.AddNotes(c.GetNoteFromCommentGroup(ctx, st.Doc, st.Comment)...)
	return res
}

func (c *CodeDecoder) GetFuncsFromASTFuncType(ctx DecoderContext, receiver gocoder.Receiver, name string, se *ast.FuncType) gocoder.Func {
	var args []gocoder.Arg
	var returns []gocoder.Arg

	if se.Params != nil {
		for _, arg := range se.Params.List {
			argType := c.GetTypeFromASTNode(ctx, arg.Type)
			fieldName := ""
			for _, argName := range arg.Names {
				fieldName = argName.Name
			}
			gocoderArg := gocoder.NewArg(fieldName, argType, false)
			gocoderArg.AddNotes(c.GetNoteFromCommentGroup(ctx, arg.Doc, arg.Comment)...)
			args = append(args, gocoderArg)
		}
	}
	if se.Results != nil {
		for _, arg := range se.Results.List {
			argType := c.GetTypeFromASTNode(ctx, arg.Type)
			fieldName := ""
			for _, argName := range arg.Names {
				fieldName = argName.Name
			}
			gocoderArg := gocoder.NewArg(fieldName, argType, false)
			gocoderArg.AddNotes(c.GetNoteFromCommentGroup(ctx, arg.Doc, arg.Comment)...)
			returns = append(returns, gocoderArg)
		}
	}

	return gocoder.NewFunc(gocoder.FuncTypeDefault, name, receiver, args, returns)
}
