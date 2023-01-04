package gocoder

// Func type
type Func interface {
	Codable
	NoteCode

	// Getter
	GetType() FuncType
	GetName() string
	GetCodes() []Codable
	GetArgs() []Arg
	GetReturns() []Arg
	GetReturnTypes() []Type
	GetReceiver() Receiver

	C(...Codable) Func
	Call(...interface{}) Value
	ToCode() Code

	InterfaceForFunc() bool
}

// FuncType type
type FuncType int

// FuncType type
const (
	FuncTypeDefault FuncType = 0
	FuncTypeInline  FuncType = 1
)

var _ Func = (*tFunc)(nil)

type tFunc struct {
	TNoteCode
	Type     FuncType
	Name     string
	Codes    []Codable
	Args     []Arg
	Returns  []Arg
	Receiver Receiver
}

func (t *tFunc) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tFunc) GetType() FuncType {
	return t.Type
}

func (t *tFunc) GetName() string {
	return t.Name
}

func (t *tFunc) GetCodes() []Codable {
	return t.Codes
}

func (t *tFunc) GetArgs() []Arg {
	return t.Args
}

func (t *tFunc) GetReturns() []Arg {
	return t.Returns
}

func (t *tFunc) GetReceiver() Receiver {
	return t.Receiver
}

func (t *tFunc) ToCode() Code {
	return &tFuncCode{
		tFunc: t,
	}
}

func (t *tFunc) InterfaceForFunc() bool {
	return true
}

func (t *tFunc) C(cs ...Codable) Func {
	t.Codes = append(t.Codes, cs...)
	return t
}

func (t *tFunc) GetReturnTypes() []Type {
	res := make([]Type, 0, len(t.Returns))
	for _, r := range t.Returns {
		res = append(res, r.GetType())
	}
	return res
}

func (t *tFunc) Call(argsI ...interface{}) Value {
	args := MustToValueList(argsI...)
	var retType Type
	if len(t.Returns) > 0 {
		retType = t.Returns[0].GetType()
	}
	return &tValue{
		Action:       ValueActionFuncCall,
		IType:        retType,
		CallReturns:  t.GetReturnTypes(),
		CallArgs:     args,
		Func:         t,
		Left:         nil,
		Right:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Values:       nil,
		CallArgTypes: nil,
	}
}

type tFuncCode struct {
	*tFunc
}

func (t *tFuncCode) C(cs ...Codable) Code {
	t.Codes = append(t.Codes, cs...)
	return t
}
