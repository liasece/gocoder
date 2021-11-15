package gocoder

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// Codeable type
type Codeable interface {
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
	if i, ok := i.(Type); ok {
		return i
	}
	if i, ok := i.(Value); ok {
		if it := i.GetIType(); it != nil && !it.IsNil() {
			return it
		}
		return i.Type()
	}
	if i, ok := i.(reflect.Type); ok {
		return NewType(i)
	}
	return NewTypeI(i)
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
func NewForRange(autoSet bool, typ FuncType, toValues Value, value Value, cs ...Codeable) ForRange {
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
		Name:  name,
		IType: t,
	}
}

// NewValueFunc func
func NewValueFunc(name string, typ Type, argTypes []Type, returns []Type) Value {
	return &tValue{
		Name:         name,
		IType:        typ,
		CallArgTypes: argTypes,
		CallReturns:  returns,
	}
}

// NewValueNameI func
func NewValueNameI(name string, i interface{}) Value {
	return &tValue{
		Name:   name,
		IType:  NewType(reflect.TypeOf(i)),
		IValue: i,
	}
}

// NewOnlyTypeValue func
func NewOnlyTypeValue(t Type) Value {
	return &tValue{
		IType: t,
	}
}

// NewValues func
func NewValues(vs ...Value) Value {
	return &tValue{
		Values: vs,
	}
}

// NewValueNameRef func
func NewValueNameRef(name string, t reflect.Type) Value {
	return &tValue{
		Name:  name,
		IType: NewType(t),
	}
}

// NewValueI func
func NewValueI(i interface{}) Value {
	return &tValue{
		IType:  NewType(reflect.TypeOf(i)),
		IValue: i,
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
		Type: reflect.TypeOf(i),
	}
}

// NewTypeName func
func NewTypeName(name string) Type {
	return &tType{
		Str: name,
	}
}

// NewTypeDetail func
func NewTypeDetail(pkg string, name string) Type {
	return &tType{
		Str: name,
		Pkg: pkg,
	}
}

// NewIf func
func NewIf(v Value, cs ...Codeable) If {
	return &tIf{
		IfV:   v,
		Codes: cs,
	}
}

// NewArgI func
func NewArgI(name string, i interface{}) Arg {
	return &tArg{
		Name: name,
		Type: NewTypeI(i),
	}
}

// NewArg func
func NewArg(name string, typ Type) Arg {
	return &tArg{
		Name: name,
		Type: typ,
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
func NewFunc(typ FuncType, name string, receiver Receiver, args []Arg, returns []Type, notes ...Note) Func {
	return &tFunc{
		Type:     typ,
		Name:     name,
		Receiver: receiver,
		Args:     args,
		Returns:  returns,
		Notes:    notes,
	}
}

// NewStruct func
func NewStruct(name string, fs []Field) Struct {
	return &tStruct{
		ReName: name,
		Fields: fs,
	}
}

// NewField func
func NewField(name string, typ Type, tag string) Field {
	return &tField{
		Type:   typ,
		ReName: name,
		Tag:    tag,
	}
}

// NewType func
func NewType(t reflect.Type) Type {
	return &tType{
		Type: t,
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
	return &tCode{}
}

// code type
type Code interface {
	Codeable
	C(cs ...Codeable) Code
	GetCodes() []Codeable
}

type tCode struct {
	Codes []Codeable
}

func (t *tCode) GetCodes() []Codeable {
	return t.Codes
}

func (t *tCode) WriteCode(w Writer) {
	w.WriteCode(t)
}

// C func
func (t *tCode) C(cs ...Codeable) Code {
	t.Codes = append(t.Codes, cs...)
	return t
}
