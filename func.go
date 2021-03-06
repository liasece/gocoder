package gocoder

// Func type
type Func interface {
	Codeable

	// Getter
	GetType() FuncType
	GetName() string
	GetCodes() []Codeable
	GetArgs() []Arg
	GetReturns() []Type
	GetNotes() []Note
	GetReceiver() Type

	C(...Codeable) Func
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

type tFunc struct {
	Type     FuncType
	Name     string
	Codes    []Codeable
	Args     []Arg
	Returns  []Type
	Receiver Receiver
	Notes    []Note
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

func (t *tFunc) GetCodes() []Codeable {
	return t.Codes
}

func (t *tFunc) GetArgs() []Arg {
	return t.Args
}

func (t *tFunc) GetReturns() []Type {
	return t.Returns
}

func (t *tFunc) GetReceiver() Type {
	return t.Receiver
}

func (t *tFunc) GetNotes() []Note {
	return t.Notes
}

func (t *tFunc) ToCode() Code {
	return &tFuncCode{
		tFunc: t,
	}
}

func (t *tFunc) InterfaceForFunc() bool {
	return true
}

func (t *tFunc) C(cs ...Codeable) Func {
	t.Codes = append(t.Codes, cs...)
	return t
}

func (t *tFunc) Call(argsI ...interface{}) Value {
	args := MustToValueList(argsI...)
	var retType Type
	if len(t.Returns) > 0 {
		retType = t.Returns[0]
	}
	return &tValue{
		Action:      ValueActionFuncCall,
		IType:       retType,
		CallReturns: t.Returns,
		CallArgs:    args,
		Func:        t,
	}
}

type tFuncCode struct {
	*tFunc
}

func (t *tFuncCode) C(cs ...Codeable) Code {
	t.Codes = append(t.Codes, cs...)
	return t
}
