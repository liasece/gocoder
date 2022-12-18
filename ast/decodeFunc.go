package ast

import (
	"go/ast"

	"github.com/liasece/gocoder"
)

func (c *ASTCoder) GetFuncsFromASTFieldList(receiver gocoder.Receiver, st *ast.FieldList, opt *gocoder.ToCodeOption) ([]gocoder.Func, error) {
	fs := make([]gocoder.Func, 0)
	for _, arg := range st.List {
		f, err := c.GetFuncsFromASTField(receiver, arg, opt)
		if err != nil {
			// return nil, err
		} else {
			fs = append(fs, f)
		}
	}
	return fs, nil
}

func (c *ASTCoder) GetFuncsFromASTFuncDecl(st *ast.FuncDecl, opt *gocoder.ToCodeOption) (gocoder.Func, error) {
	// ast.Print(c.fset, st)
	var name string
	{
		// get name
		name = st.Name.Name
	}
	receiver, err := c.GetReceiverFromASTField(st.Recv.List[0], opt)
	if err != nil {
		return nil, err
	}
	return c.GetFuncsFromASTFuncType(receiver, name, st.Type, opt)
}

func (c *ASTCoder) GetFuncsFromASTField(receiver gocoder.Receiver, st *ast.Field, opt *gocoder.ToCodeOption) (gocoder.Func, error) {
	// ast.Print(c.fset, st)
	var name string
	{
		// get name
		if len(st.Names) > 0 {
			name = st.Names[0].Name
		}
	}
	return c.GetFuncsFromASTFuncType(receiver, name, st.Type.(*ast.FuncType), opt)
}

func (c *ASTCoder) GetFuncsFromASTFuncType(receiver gocoder.Receiver, name string, se *ast.FuncType, opt *gocoder.ToCodeOption) (gocoder.Func, error) {
	var args []gocoder.Arg
	var returns []gocoder.Type

	if se.Params != nil {
		for _, arg := range se.Params.List {
			argType, err := c.GetTypeFromASTExpr(arg.Type, opt)
			if err != nil {
				return nil, err
			}
			for _, argName := range arg.Names {
				args = append(args, gocoder.NewArg(argName.Name, argType, false))
			}
		}
	}
	if se.Results != nil {
		for _, arg := range se.Results.List {
			argType, err := c.GetTypeFromASTExpr(arg.Type, opt)
			if err != nil {
				return nil, err
			}
			returns = append(returns, argType)
		}
	}

	return gocoder.NewFunc(gocoder.FuncTypeDefault, name, receiver, args, returns), nil
}
