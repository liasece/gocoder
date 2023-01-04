package gocoder

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// Codable type
type Codable interface {
	WriteCode(Writer)
}

// MustToNoteList func
func MustToNoteList(kind NoteKind, is ...interface{}) []Note {
	res := make([]Note, 0, len(is))
	for _, v := range is {
		res = append(res, MustToNote(kind, v))
	}
	return res
}

// MustToNote func
func MustToNote(kind NoteKind, i interface{}) Note {
	switch i := i.(type) {
	case Note:
		return i
	case string:
		return NewNote(i, kind)
	default:
		panic(fmt.Sprint("MustToNote unknown type", reflect.TypeOf(i).String()))
	}
}

// MustToType func
func MustToType(i interface{}) Type {
	switch i := i.(type) {
	case Type:
		return i
	case Value:
		if it := i.GetIType(); it != nil && !it.IsNil() {
			return it
		}
		return i.Type()
	case reflect.Type:
		return NewType(i)
	default:
		return NewTypeI(i)
	}
}

// MustToValueList func
func MustToValueList(is ...interface{}) []Value {
	vs := make([]Value, 0, len(is))
	for _, i := range is {
		switch i := i.(type) {
		case []Value:
			vs = append(vs, i...)
		default:
			vs = append(vs, MustToValue("", i))
		}
	}
	return vs
}

// MustToValues func
func MustToValues(is ...interface{}) Value {
	return NewValues(MustToValueList(is...)...)
}

// MustToValue func
func MustToValue(name string, i interface{}) Value {
	if i == nil {
		return NewValueNameI(name, nil)
	}
	if i, ok := i.(Value); ok {
		return i
	}
	if i, ok := i.([]Value); ok {
		is := make([]interface{}, 0, len(i))
		for _, v := range i {
			is = append(is, v)
		}
		return MustToValues(is...)
	}
	if i, ok := i.(Type); ok {
		return NewValue(name, i)
	}
	if i, ok := i.(reflect.Type); ok {
		return NewValueNameRef(name, i)
	}
	// check is package function
	{
		value := reflect.ValueOf(i)
		if value.Kind() == reflect.Func {
			return pkgFunc(i)
		}
	}
	if name == "" {
		return NewValueI(i)
	}
	return NewValueNameI(name, i)
}

// pkgFunc func
func pkgFunc(funcI interface{}) Value {
	value := reflect.ValueOf(funcI)
	typ := value.Type()
	if value.Kind() != reflect.Func {
		panic("pkgFunc funcI not a func")
	}
	pc := runtime.FuncForPC(value.Pointer())
	pcName := pc.Name()
	pkg := ""
	funcName := pcName
	i := strings.LastIndex(pcName, ".")
	if i >= 0 {
		pkg = pcName[:i]
		funcName = pcName[i:]
	}
	_, name := filepath.Split(pcName)
	valueName := pcName
	if name != "" {
		valueName = name
	}
	ins := make([]Type, 0)
	outs := make([]Type, 0)
	for i := 0; i < typ.NumIn(); i++ {
		ins = append(ins, NewType(typ.In(i)))
	}
	for i := 0; i < typ.NumOut(); i++ {
		outs = append(outs, NewType(typ.Out(i)))
	}
	return NewValueFunc(valueName, NewTypeDetail(pkg, funcName), ins, outs)
}

// NewForRange func
func NewForRange(autoSet bool, typ FuncType, toValues Value, value Value, cs ...Codable) ForRange {
	return &tForRange{
		Type:     typ,
		AutoSet:  autoSet,
		ToValues: toValues,
		Value:    value,
		Codes:    cs,
	}
}

// NewPtrChecker func
func NewPtrChecker(ifNotNil bool, checkerValue ...Value) PtrChecker {
	return &tPtrChecker{
		CheckerValue: checkerValue,
		IfNotNil:     ifNotNil,
		Handlers:     nil,
	}
}

// NewNote func
func NewNote(content string, kind NoteKind) Note {
	return &tNote{
		Content: content,
		Kind:    kind,
	}
}

// NewValue func
func NewValue(name string, t Type) Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Name:         name,
		IType:        t,
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

// NewValueFunc func
func NewValueFunc(name string, typ Type, argTypes []Type, returns []Type) Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Name:         name,
		IType:        typ,
		CallArgTypes: argTypes,
		CallReturns:  returns,
		Left:         nil,
		Action:       "",
		Right:        nil,
		IValue:       nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
	}
}

// NewValueNameI func
func NewValueNameI(name string, i interface{}) Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Name:         name,
		IType:        NewType(reflect.TypeOf(i)),
		IValue:       i,
		Left:         nil,
		Action:       "",
		Right:        nil,
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

// NewOnlyTypeValue func
func NewOnlyTypeValue(t Type) Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		IType:        t,
		Left:         nil,
		Action:       "",
		Right:        nil,
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

// NewValues func
func NewValues(vs ...Value) Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Values:       vs,
		Left:         nil,
		Action:       "",
		Right:        nil,
		Name:         "",
		IValue:       nil,
		Str:          "",
		Func:         nil,
		IType:        nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

// NewValueNameRef func
func NewValueNameRef(name string, t reflect.Type) Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		Name:         name,
		IType:        NewType(t),
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

// NewValueI func
func NewValueI(i interface{}) Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		IType:        NewType(reflect.TypeOf(i)),
		IValue:       i,
		Left:         nil,
		Action:       "",
		Right:        nil,
		Name:         "",
		Str:          "",
		Func:         nil,
		Values:       nil,
		CallArgs:     nil,
		CallArgTypes: nil,
		CallReturns:  nil,
	}
}

// NewValueNil func
func NewValueNil() Value {
	return NewValueNameI("nil", nil)
}

// NewValueNone func
func NewValueNone() Value {
	return NewValueNameI("_", nil)
}

// NewTypeI func
func NewTypeI(i interface{}) Type {
	return &tType{
		TNoteCode:   TNoteCode{nil},
		Type:        reflect.TypeOf(i),
		Str:         "",
		Pkg:         "",
		Named:       "",
		Next:        nil,
		inReference: false,
		kind:        0,
		funcs:       nil,
		fields:      nil,
	}
}

// NewTypeName func
func NewTypeName(name string) Type {
	return NewTypeDetail("", name)
}

// NewTypeDetail func
func NewTypeDetail(pkg string, name string) Type {
	var next Type
	if strings.HasPrefix(name, "*") && name != "*" {
		next = NewTypeDetail(pkg, name[1:])
		name = name[:1]
	} else if strings.HasPrefix(name, "[]") && name != "[]" {
		next = NewTypeDetail(pkg, name[2:])
		name = name[:2]
	}
	return &tType{
		TNoteCode:   TNoteCode{nil},
		Str:         name,
		Pkg:         pkg,
		Type:        nil,
		Named:       "",
		Next:        next,
		inReference: false,
		kind:        0,
		funcs:       nil,
		fields:      nil,
	}
}

// NewIf func
func NewIf(v Value, cs ...Codable) If {
	return &tIf{
		IfV:   v,
		Codes: cs,
		IPre:  nil,
		INext: nil,
	}
}

// NewArgI func
func NewArgI(name string, i interface{}) Arg {
	return &tArg{
		TNoteCode:      TNoteCode{nil},
		Name:           name,
		Type:           NewTypeI(i),
		VariableLength: false,
	}
}

// NewArg func
func NewArg(name string, typ Type, variableLength bool) Arg {
	return &tArg{
		TNoteCode:      TNoteCode{nil},
		Name:           name,
		Type:           typ,
		VariableLength: variableLength,
	}
}

// NewReceiver func
func NewReceiver(name string, typ Type) Receiver {
	return &tReceiver{
		Type:   typ,
		ReName: name,
	}
}

// NewFunc func
func NewFunc(typ FuncType, name string, receiver Receiver, args []Arg, returns []Arg, notes ...Note) Func {
	f := &tFunc{
		TNoteCode: TNoteCode{nil},
		Type:      typ,
		Name:      name,
		Receiver:  receiver,
		Args:      args,
		Returns:   returns,
		Codes:     nil,
	}
	f.SetNotes(notes)
	return f
}

// NewStruct func
func NewStruct(name string, fs []Field) Type {
	return &tType{
		TNoteCode:   TNoteCode{nil},
		Named:       name,
		fields:      fs,
		inReference: false,
		Str:         "",
		Pkg:         "",
		Type:        nil,
		Next:        nil,
		kind:        reflect.Struct,
		funcs:       nil,
	}
}

// NewInterface func
func NewInterface(name string, fs []Func) Type {
	return &tType{
		TNoteCode:   TNoteCode{nil},
		Named:       name,
		funcs:       fs,
		inReference: false,
		Str:         "",
		Pkg:         "",
		Type:        nil,
		Next:        nil,
		kind:        reflect.Interface,
		fields:      nil,
	}
}

// NewField func
func NewField(name string, typ Type, tag string) Field {
	return &tField{
		TNoteCode: TNoteCode{nil},
		Type:      typ,
		ReName:    name,
		Tag:       tag,
	}
}

// NewType func
func NewType(t reflect.Type) Type {
	return &tType{
		TNoteCode:   TNoteCode{nil},
		Type:        t,
		Str:         "",
		Pkg:         "",
		Named:       "",
		Next:        nil,
		inReference: false,
		kind:        0,
		funcs:       nil,
		fields:      nil,
	}
}

// NewReturn func
func NewReturn(v Value) Return {
	return &tReturn{
		Value: v,
	}
}

// NewCode func
func NewCode() Code {
	return &tCode{
		Codes: nil,
	}
}

// code type
type Code interface {
	Codable
	C(cs ...Codable) Code
	GetCodes() []Codable
}

type tCode struct {
	Codes []Codable
}

func (t *tCode) GetCodes() []Codable {
	return t.Codes
}

func (t *tCode) WriteCode(w Writer) {
	w.WriteCode(t)
}

// C func
func (t *tCode) C(cs ...Codable) Code {
	t.Codes = append(t.Codes, cs...)
	return t
}
