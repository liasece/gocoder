package gocoder

import (
	"fmt"
	"reflect"
	"strings"
)

// Value type
type Value interface {
	Codable
	NoteCode

	// Getter
	GetAction() ValueAction
	GetName() string
	GetSrcValue() interface{}
	GetLeft() Codable
	GetRight() Value
	GetFunc() Func
	GetCallArgs() []Value
	GetValues() []Value
	GetIType() Type

	ToArg() Arg
	Call(...interface{}) Value
	TypeString() string
	Field(name string) Value
	Method(method string) Value
	Note(...interface{}) Value
	ToPtr() Value
	Type() Type
	IsNilType() bool
	IsPtr() bool
	Depth() int
	NeedParent() bool
	Returns() []Value
	Cast(interface{}) Value
	Assertion(interface{}) Value
	Dot(name string) Value
	Set(interface{}, ...*SetOption) Value
	AutoSet(interface{}, ...*SetOption) Value
	Check() Value
	Index(interface{}) Value
	Add(interface{}) Value
	Sub(interface{}) Value
	Mul(interface{}) Value
	Div(interface{}) Value
	Equal(interface{}) Value
	GT(interface{}) Value
	LT(interface{}) Value
	GE(interface{}) Value
	LE(interface{}) Value
	NE(interface{}) Value
	Not() Value
	Or(interface{}) Value
	And(interface{}) Value
	UnPtr() Value
	TakePtr() Value
}

// ValueAction type
type ValueAction string

// ValueAction type
const (
	ValueActionNone          ValueAction = ""
	ValueActionDot           ValueAction = "."
	ValueActionSet           ValueAction = "="
	ValueActionAutoSet       ValueAction = ":="
	ValueActionAdd           ValueAction = "+"
	ValueActionSub           ValueAction = "-"
	ValueActionMul           ValueAction = "*"
	ValueActionDiv           ValueAction = "/"
	ValueActionEqual         ValueAction = "=="
	ValueActionGT            ValueAction = ">"
	ValueActionLT            ValueAction = "<"
	ValueActionGE            ValueAction = ">="
	ValueActionLE            ValueAction = "<="
	ValueActionNE            ValueAction = "!="
	ValueActionOr            ValueAction = "||"
	ValueActionAnd           ValueAction = "&&"
	ValueActionNot           ValueAction = "!"
	ValueActionUnPtr         ValueAction = "**"
	ValueActionTakePtr       ValueAction = "&"
	ValueActionCastType      ValueAction = "()()"
	ValueActionAssertionType ValueAction = ".()"
	ValueActionFuncCall      ValueAction = "()"
	ValueActionIndex         ValueAction = "[]"
	ValueActionZero          ValueAction = "0"
)

var actionCodeConv = map[ValueAction]string{
	ValueActionUnPtr: "*",
}

var valueActionNeedParent = map[ValueAction]bool{
	ValueActionUnPtr: true,
	ValueActionAdd:   true,
	ValueActionSub:   true,
	ValueActionMul:   true,
	ValueActionDiv:   true,
	// ValueActionEqual:   true,
	// ValueActionGT:      true,
	// ValueActionLT:      true,
	// ValueActionGE:      true,
	// ValueActionLE:      true,
	// ValueActionNE:      true,
	ValueActionOr:      true,
	ValueActionAnd:     true,
	ValueActionNot:     true,
	ValueActionTakePtr: true,
}

var valueActionForceNeedParent = map[ValueAction]bool{
	ValueActionUnPtr: true,
}

type tValue struct {
	TNoteCode

	Left   Codable
	Action ValueAction
	Right  Value

	Name   string
	IType  Type
	IValue interface{}
	Str    string
	Func   Func

	Values []Value

	CallArgs     []Value
	CallArgTypes []Type
	CallReturns  []Type
}

func (t *tValue) GetAction() ValueAction {
	return t.Action
}

func (t *tValue) GetName() string {
	return t.Name
}

func (t *tValue) GetSrcValue() interface{} {
	return t.IValue
}

func (t *tValue) GetLeft() Codable {
	return t.Left
}

func (t *tValue) GetRight() Value {
	return t.Right
}

func (t *tValue) GetCallArgs() []Value {
	return t.CallArgs
}

func (t *tValue) GetArgs() []Value {
	return t.CallArgs
}

func (t *tValue) GetFunc() Func {
	return t.Func
}

func (t *tValue) GetValue() Value {
	return t
}

func (t *tValue) GetIType() Type {
	return t.IType
}

func (t *tValue) GetValues() []Value {
	return t.Values
}

func (t *tValue) GetReturns() []Value {
	res := make([]Value, len(t.CallReturns))
	for i, v := range t.CallReturns {
		res[i] = NewValue("", v)
	}
	return res
}

func (t *tValue) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tValue) Field(name string) Value {
	return t.field(name)
}

func (t *tValue) field(name string) *tValue {
	names := strings.Split(name, ".")
	if len(names) > 1 {
		t = t.field(strings.Join(names[:len(names)-1], "."))
		name = names[len(names)-1]
	}
	noPtrT := t.Type()
	if noPtrT.Kind() == reflect.Ptr {
		noPtrT = noPtrT.Elem()
	}
	if noPtrT.Kind() != reflect.Struct {
		panic("value isn't struct")
	}
	f := noPtrT.FieldByName(name)
	if f == nil {
		panic(fmt.Sprintf("value no target field: %s", name))
	}
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Name:         t.Name + "." + f.GetName(),
		IType:        f.GetType(),
		Left:         nil,
		Action:       "",
		Right:        nil,
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Method(name string) Value {
	noPtrT := t.Type()
	if t.CallReturns != nil {
		if len(t.CallReturns) > 1 {
			panic("len(t.CallReturns) > 1")
		}
		noPtrT = t.CallReturns[0]
	}
	if noPtrT != nil {
		// if noPtrT.Kind() == reflect.Ptr {
		// 	noPtrT = noPtrT.Elem()
		// }
		if noPtrT.RefType() == nil {
			return &tValue{
				TNoteCode:    TNoteCode{nil},
				Left:         t,
				Name:         name,
				IType:        nil,
				Action:       "",
				Right:        nil,
				IValue:       nil,
				Str:          "",
				Func:         nil,
				Values:       nil,
				CallArgs:     nil,
				CallArgTypes: nil,
				CallReturns:  nil,
			}
		}
		if f, ok := noPtrT.MethodByName(name); !ok {
			panic(fmt.Sprintf("value isn't struct, type: %s(%s) name: %s", noPtrT.String(), noPtrT.RefType().String(), name))
		} else {
			outTypes := make([]Type, f.Type.NumOut())
			for i := range outTypes {
				outTypes[i] = NewType(f.Type.Out(i))
			}
			inTypes := make([]Type, f.Type.NumIn()-1)
			for i := range inTypes {
				inTypes[i] = NewType(f.Type.In(i + 1))
			}
			return &tValue{
				TNoteCode:    TNoteCode{nil},
				Left:         t,
				Name:         f.Name,
				IType:        NewType(f.Type),
				CallReturns:  outTypes,
				CallArgTypes: inTypes,
				Action:       "",
				Right:        nil,
				IValue:       nil,
				Str:          "",
				Func:         nil,
				Values:       nil,
				CallArgs:     nil,
			}
		}
	}
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Name:         name,
		IType:        nil,
		Action:       "",
		Right:        nil,
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) ToArg() Arg {
	return &tArg{
		TNoteCode:      t.TNoteCode.Clone(),
		Name:           t.Name,
		Type:           t.Type(),
		VariableLength: false,
	}
}

func (t *tValue) Cast(i interface{}) Value {
	target := MustToType(i)
	if t.IsNilType() || target == nil || target.IsNil() {
		return t
	}
	var midValue Value = t
	if target.IsPtr() && !midValue.IsPtr() {
		midValue = midValue.TakePtr()
	} else if !target.IsPtr() && midValue.IsPtr() {
		midValue = midValue.UnPtr()
	}
	if t.Type() != nil {
		noPtrTarget := target.UnPtr()
		noPtrCurent := t.Type().UnPtr()
		if t.ReturnType1() != nil {
			noPtrCurent = t.ReturnType1().UnPtr()
		}
		if noPtrTarget.Kind() == reflect.Interface {
			if !noPtrCurent.Implements(noPtrTarget.RefType()) {
				panic("can't (" + noPtrCurent.String() + ") Implements (" + noPtrTarget.String() + ")")
			}
		} else if noPtrCurent.Kind() != noPtrTarget.Kind() && !noPtrCurent.IsStruct() && !noPtrTarget.IsStruct() {
			if !noPtrCurent.ConvertibleTo(noPtrTarget.RefType()) || !noPtrTarget.ConvertibleTo(noPtrCurent.RefType()) {
				panic("can't (" + noPtrCurent.String() + ") ConvertibleTo (" + noPtrTarget.String() + ")")
			} else {
				midValue = &tValue{
					TNoteCode:    TNoteCode{nil},
					Left:         target,
					Action:       ValueActionCastType,
					Right:        midValue,
					IType:        target,
					Name:         "",
					IValue:       nil,
					Str:          "",
					Func:         nil,
					Values:       nil,
					CallArgs:     nil,
					CallArgTypes: nil,
					CallReturns:  nil,
				}
			}
		}
		if noPtrCurent.Kind() == noPtrTarget.Kind() && noPtrCurent.String() != noPtrTarget.String() {
			midValue = &tValue{
				TNoteCode:    TNoteCode{nil},
				Left:         target,
				Action:       ValueActionCastType,
				Right:        midValue,
				IType:        target,
				Name:         "",
				IValue:       nil,
				Str:          "",
				Func:         nil,
				Values:       nil,
				CallArgs:     nil,
				CallArgTypes: nil,
				CallReturns:  nil,
			}
		}
	} else {
		panic(fmt.Sprintf("Cast error %v", t.TypeString()))
	}
	return midValue
}

func (t *tValue) Assertion(i interface{}) Value {
	target := MustToType(i)
	if t.IsNilType() || target == nil || target.IsNil() {
		if target != nil && target.String() != "" {
			return &tValue{
				TNoteCode:    TNoteCode{nil},
				Left:         t,
				Action:       ValueActionAssertionType,
				Right:        t,
				IType:        target,
				Name:         "",
				IValue:       nil,
				Str:          "",
				Func:         nil,
				Values:       nil,
				CallArgs:     nil,
				CallArgTypes: nil,
				CallReturns:  nil,
			}
		}
		return t
	}
	var midValue Value = t
	if t.Type() != nil {
		noPtrTarget := target.UnPtr()
		noPtrCurent := t.Type().UnPtr()
		if t.ReturnType1() != nil {
			noPtrCurent = t.ReturnType1().UnPtr()
		}
		if noPtrCurent.Kind() == reflect.Interface {
			if !noPtrTarget.Implements(noPtrCurent.RefType()) {
				panic("can't (" + noPtrTarget.String() + ") Implements (" + noPtrCurent.String() + ")")
			} else {
				return &tValue{
					TNoteCode:    TNoteCode{nil},
					Left:         t,
					Action:       ValueActionAssertionType,
					Right:        midValue,
					IType:        target,
					Name:         "",
					IValue:       nil,
					Str:          "",
					Func:         nil,
					Values:       nil,
					CallArgs:     nil,
					CallArgTypes: nil,
					CallReturns:  nil,
				}
			}
		}
	}
	panic(fmt.Sprintf("Assertion error %v to %v", t.TypeString(), target.String()))
}

func (t *tValue) Call(argsI ...interface{}) Value {
	args := MustToValueList(argsI...)
	if t.CallArgTypes != nil {
		if len(t.CallArgTypes) > 0 && t.CallArgTypes[len(t.CallArgTypes)-1].Kind() == reflect.Slice {
			lastArg := t.CallArgTypes[len(t.CallArgTypes)-1].Elem()
			if len(args) < len(t.CallArgTypes)-1 {
				panic(fmt.Sprintf("len(args)(%d)(%+v) != len(t.CallArgTypes)-1(%d)(%+v), func: %s %s", len(args), args, len(t.CallArgTypes), t.CallArgTypes, t.Type().Name(), t.Type().String()))
			}
			for i := len(t.CallArgTypes) - 1; i < len(args); i++ {
				if args[i].Type().RefType() != lastArg.RefType() && (args[i].Type().RefType() != nil && !args[i].Type().RefType().Implements(lastArg.RefType())) {
					panic(fmt.Sprintf("args[i].Type().RefType()(%v) != lastArg.RefType()(%v), func: %s %s\n %v : %v", args[i].TypeString(), lastArg.String(), t.Type().Name(), t.Type().String(), args, t.CallArgTypes))
				}
			}
		} else if len(args) != len(t.CallArgTypes) {
			panic(fmt.Sprintf("len(args)(%d)(%+v) != len(t.CallArgTypes)(%d)(%+v), func: %s %s", len(args), args, len(t.CallArgTypes), t.CallArgTypes, t.Type().Name(), t.Type().String()))
		}
	}
	realArgs := args
	if len(t.CallArgTypes) > 0 {
		if t.CallArgTypes[len(t.CallArgTypes)-1].Kind() != reflect.Slice {
			for i, arg := range realArgs {
				if i >= len(t.CallArgTypes) {
					realArgs[i] = arg.Cast(t.CallArgTypes[len(t.CallArgTypes)-1])
				} else {
					realArgs[i] = arg.Cast(t.CallArgTypes[i])
				}
			}
		}
	}
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionFuncCall,
		CallArgs:     realArgs,
		IType:        t.Type(),
		CallReturns:  t.CallReturns,
		Name:         "",
		Right:        nil,
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgTypes: nil,
	}
}

func (t *tValue) IsPtr() bool {
	if t.IType != nil {
		return t.IType.IsPtr()
	}
	if t.Type() == nil {
		return false
	}
	return t.Type().IsPtr()
}

func (t *tValue) Type() Type {
	if t.IType != nil {
		return t.IType
	}
	if len(t.CallReturns) > 0 {
		return t.CallReturns[0]
	}
	return nil
}

func (t *tValue) ReturnType1() Type {
	if len(t.CallReturns) > 0 {
		return t.CallReturns[0]
	}
	if t.IType != nil {
		return t.IType
	}
	return nil
}

func (t *tValue) IsNilType() bool {
	return t.Type() == nil || t.Type().IsNil()
}

func (t *tValue) TypeString() string {
	if t.Type() == nil {
		return ""
	}
	return t.Type().String()
}

func (t *tValue) Note(is ...interface{}) Value {
	t.AddNotes(MustToNoteList(NoteKindBlock, is...)...)
	return t
}

func (t *tValue) Depth() int {
	i := 1
	if t.Action == ValueActionNone {
		return i
	}
	max := 0
	if t.Left != nil {
		if lv, ok := t.Left.(Value); ok {
			tmp := lv.Depth()
			if tmp > max {
				max = tmp
			}
		}
	}
	if t.Right != nil {
		tmp := t.Right.Depth()
		if tmp > max {
			max = tmp
		}
	}
	return i + max
}

func (t *tValue) NeedParent() bool {
	if valueActionForceNeedParent[t.Action] {
		return true
	}
	if !valueActionNeedParent[t.Action] {
		return false
	}
	depth := t.Depth()
	if depth <= 1 {
		return false
	}
	if t.Left == nil && depth == 2 {
		return false
	}
	return true
}

func (t *tValue) Set(i interface{}, opts ...*SetOption) Value {
	v := MustToValue("", i)
	opt := MergeSetOpt(opts...)
	right := v
	if opt.notCast == nil || !(*opt.notCast) {
		right = right.Cast(t.Type())
	}
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionSet,
		Right:        right,
		IType:        t.Type(),
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Dot(name string) Value {
	field, ok := t.Type().FieldTypeByName(name)
	if !ok {
		panic("Value Dot FieldByName " + name + " not ok")
	}
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionDot,
		Name:         name,
		IType:        field,
		Right:        nil,
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) AutoSet(i interface{}, opts ...*SetOption) Value {
	v := MustToValue("", i)
	opt := MergeSetOpt(opts...)
	right := v
	if opt.notCast == nil || !(*opt.notCast) {
		right = right.Cast(t.Type())
	}
	if !right.IsNilType() {
		t.IType = right.Type()
	}
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionAutoSet,
		Right:        right,
		IType:        right.Type(),
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Index(i interface{}) Value {
	v := MustToValue("", i)
	left := t.UnPtr()
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         left,
		Action:       ValueActionIndex,
		Right:        v,
		IType:        left.Type().Elem(),
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Check() Value {
	if t.Type() == nil || t.Type().IsNil() {
		panic(fmt.Sprintf("tValue %+v type can't be nil", t))
	}
	return t
}

func (t *tValue) Returns() []Value {
	res := t.Values
	if t.Func != nil {
		res = append(res, t.GetReturns()...)
	}
	return res
}

func (t *tValue) Add(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionAdd,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Sub(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionSub,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Mul(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionMul,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Div(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionDiv,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Equal(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionEqual,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) GT(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionGT,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) LT(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionLT,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) GE(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionGE,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) LE(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionLE,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) NE(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionNE,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Not() Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Action:       ValueActionNot,
		Right:        t,
		IType:        nil,
		Name:         "",
		Left:         nil,
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) Or(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionOr,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) And(i interface{}) Value {
	v := MustToValue("", i)
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Left:         t,
		Action:       ValueActionAnd,
		Right:        v,
		IType:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) UnPtr() Value {
	if t.Type() != nil && t.Type().RefType() != nil {
		if t.Type().Kind() != reflect.Ptr {
			return t
		}
	}
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Action:       ValueActionUnPtr,
		Right:        t,
		IType:        t.Type().Elem(),
		Name:         "",
		Left:         nil,
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) TakePtr() Value {
	if !t.IsNilType() {
		if t.Type().Kind() == reflect.Ptr {
			return t
		}
	}
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Action:       ValueActionTakePtr,
		Right:        t,
		IType:        t.Type().TackPtr(),
		Name:         "",
		Left:         nil,
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

func (t *tValue) ToPtr() Value {
	tmp := *t
	res := &tmp
	res.IType = res.Type().TackPtr()
	return res
}

func (t *tValue) ToNoPtr() Value {
	tmp := *t
	res := &tmp
	res.IType = res.Type().UnPtr()
	return res
}
