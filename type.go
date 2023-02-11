package gocoder

import (
	"reflect"
	"strings"
)

// Type type
type Type interface {
	Codable
	NoteCode

	RefType() reflect.Type
	IsPtr() bool
	IsSlice() bool
	IsMap() bool
	UnPtr() Type
	IsStruct() bool
	TackPtr() Type
	Slice() Type
	IsNil() bool
	Elem() Type
	Kind() reflect.Kind
	SetKind(reflect.Kind)
	String() string
	ShowString() string
	CurrentCode() string
	Package() string
	PackageInReference() string // like if pkg is `github.com/liasece/gocoder`, return `gocoder`
	ConvertibleTo(i interface{}) bool
	Implements(i interface{}) bool
	NumField() int
	Field(i int) Field
	FieldTypeByName(name string) (Type, bool)
	MethodByName(name string) (reflect.Method, bool)
	Zero() Value
	Name() string
	GetNamed() string
	GetRowStr() string
	SetNamed(string)
	WarpNamed(named string) Type // build type like `type Foo SubType`, named is `Foo`
	SetPkg(string)
	GetNext() Type

	AllSub() []Type // list all type chian nodes, top type in last index
	HasPtrSubType() bool
	HasSliceSubType() bool
	InterfaceForType() bool
	InReference() bool
	SetInReference(bool)

	// struct
	GetFields() []Field
	AddFields([]Field)
	FieldByName(name string) Field

	// interface
	GetFuncs() []Func
	FuncByName(name string) Func

	Clone() Type
}

var _ Type = (*tType)(nil)

type tType struct {
	TNoteCode
	reflect.Type

	Str   string // like `*` or `[]` or `map` or `struct` or `interface` or `int` or `int64` or `string` or `Time`
	Pkg   string // like `time` or `github.com/liasece/gocoder`
	Named string // like `Foo` in `type Foo string`
	Next  Type

	// struct
	fields []Field

	// interface
	funcs []Func

	inReference bool // like: `boo int`, the int is inReference
	kind        reflect.Kind
}

func (t *tType) Clone() Type {
	res := &tType{
		TNoteCode: t.TNoteCode.Clone(),
		Type:      t.Type,

		Str:         t.Str,
		Pkg:         t.Pkg,
		Named:       t.Named,
		Next:        t.Next,
		inReference: t.inReference,
		kind:        t.kind,
		fields:      t.fields,
		funcs:       t.funcs,
	}
	if t.Next != nil {
		res.Next = t.Next.Clone()
	}
	if t.fields != nil {
		res.fields = make([]Field, len(t.fields))
		for i, v := range t.fields {
			res.fields[i] = v.Clone()
		}
	}
	return res
}

func (t *tType) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tType) IsStruct() bool {
	return t.Kind() == reflect.Struct
}

func (t *tType) IsNil() bool {
	return t == nil || t.Type == nil
}

func (t *tType) IsPtr() bool {
	if t.Type == nil && t.Str != "" {
		return t.Str[0] == '*'
	}
	return t.Kind() == reflect.Ptr
}

func (t *tType) IsSlice() bool {
	return t.Kind() == reflect.Slice
}

func (t *tType) IsMap() bool {
	return t.Kind() == reflect.Map
}

// list all type chian nodes, top type in last index
func (t *tType) AllSub() []Type {
	if t.Next != nil {
		return append(t.Next.AllSub(), t)
	}
	return []Type{t}
}

func (t *tType) HasPtrSubType() bool {
	ls := t.AllSub()
	for _, v := range ls {
		if v.IsPtr() {
			return true
		}
	}
	return false
}

func (t *tType) HasSliceSubType() bool {
	ls := t.AllSub()
	for _, v := range ls {
		if v.IsSlice() {
			return true
		}
	}
	return false
}

var InterfaceType reflect.Type

func init() {
	type T struct {
		A interface{}
	}
	InterfaceType = reflect.ValueOf(T{A: nil}).Field(0).Type()
}

func (t *tType) RefType() reflect.Type {
	if t.Type != nil {
		return t.Type
	}
	if t.kind == reflect.Struct && t.fields != nil {
		fields := make([]reflect.StructField, 0, len(t.fields))
		sizeOffset := uintptr(0)
		for _, v := range t.fields {
			typ := v.GetType().RefType()
			fields = append(fields, reflect.StructField{
				Name:      v.GetName(),
				Type:      typ,
				Tag:       reflect.StructTag(v.GetTag()),
				PkgPath:   t.Pkg,
				Offset:    sizeOffset,
				Index:     nil,
				Anonymous: false,
			})
			sizeOffset += typ.Size()
		}
		return reflect.StructOf(fields)
	}
	return InterfaceType
}

func (t *tType) UnPtr() Type {
	if t.Kind() == reflect.Ptr {
		if strings.HasPrefix(t.Str, "*") {
			if t.Str == "*" {
				return t.Next
			}
			res := t.Clone().(*tType)
			res.Str = t.Str[1:]
			return res
		}
		var refType reflect.Type
		if t.Type != nil {
			refType = t.Type.Elem()
		}
		return &tType{
			TNoteCode:   TNoteCode{nil},
			Type:        refType,
			Pkg:         t.Pkg,
			Str:         "",
			Named:       "",
			Next:        nil,
			inReference: t.inReference,
			kind:        0,
			fields:      nil,
			funcs:       nil,
		}
	}
	return t
}

func (t *tType) TackPtr() Type {
	if t.Type == nil {
		if !strings.HasPrefix(t.Str, "*") {
			return &tType{
				TNoteCode:   TNoteCode{nil},
				Str:         "*",
				Next:        t,
				Type:        nil,
				Pkg:         "",
				Named:       "",
				inReference: t.inReference,
				kind:        0,
				fields:      nil,
				funcs:       nil,
			}
		}
		return t
	}
	if t.Kind() != reflect.Ptr {
		return &tType{
			TNoteCode:   TNoteCode{nil},
			Type:        reflect.PtrTo(t.Type),
			Str:         t.Str,
			Pkg:         t.Pkg,
			Next:        nil,
			Named:       "",
			inReference: t.inReference,
			kind:        0,
			fields:      nil,
			funcs:       nil,
		}
	}
	return t
}

func (t *tType) Slice() Type {
	if t.Type == nil {
		return &tType{
			TNoteCode:   TNoteCode{nil},
			Str:         "[]",
			Next:        t,
			Type:        nil,
			Pkg:         "",
			Named:       "",
			inReference: t.inReference,
			kind:        0,
			fields:      nil,
			funcs:       nil,
		}
	}
	if t.Kind() != reflect.Ptr {
		str := t.Str
		if str != "" {
			str = "[]" + str
		}
		return &tType{
			TNoteCode:   TNoteCode{nil},
			Type:        reflect.SliceOf(t.Type),
			Str:         str,
			Pkg:         "",
			Named:       "",
			Next:        nil,
			inReference: t.inReference,
			kind:        0,
			fields:      nil,
			funcs:       nil,
		}
	}
	return t
}

func (t *tType) Elem() Type {
	if t.Str == "[]" || t.Str == "*" {
		return t.Next
	}
	res := t.Clone().(*tType)
	res.Type = t.Type.Elem()
	return res
}

func (t *tType) Package() string {
	if t.Pkg != "" {
		return t.Pkg
	}
	if t.Type != nil {
		return t.Type.PkgPath()
	}
	return ""
}

func (t *tType) PackageInReference() string {
	ss := strings.Split(t.Package(), "/")
	return ss[len(ss)-1]
}

func (t *tType) CurrentCode() string {
	if t.Str != "" {
		return t.Str
	}
	if t.Type != nil {
		if t.Type.Name() != "" {
			return t.Type.String()
		}
		switch t.Type.Kind() {
		case reflect.Array:
			return "[]"
		case reflect.Chan:
			return "chan"
		case reflect.Map:
			return "map[" + t.Type.Key().Name() + "]"
		case reflect.Pointer:
			return "*"
		case reflect.Slice:
			return "[]"
		default:
			str := t.Type.String()
			if str == "[]uint8" {
				str = "[]byte"
			}
			return str
		}
	}
	return t.Named
}

func (t *tType) String() string {
	res := ""
	if t.Str != "" {
		res = t.Str
		if t.Next != nil {
			res += t.Next.String()
		}
	}
	if t.Type != nil {
		str := t.Type.String()
		if str == "[]uint8" {
			str = "[]byte"
		}
		res += str
	}
	if res == "" {
		if t.Named != "" {
			res = t.Named
		} else if t.Next != nil {
			res = t.Next.String()
		}
	}
	return res
}

func (t *tType) ShowString() string {
	head := ""
	if t.Package() != "" {
		head = t.Package() + "."
	}
	if t.Named != "" {
		return head + t.Named
	}
	if t.Str != "" {
		res := head + t.Str
		if t.Next != nil {
			res += t.Next.ShowString()
		}
		return res
	}
	if t.Type != nil {
		str := t.Type.String()
		if str == "[]uint8" {
			str = "[]byte"
		}
		return str
	}

	res := head + t.Named
	if t.Next != nil {
		res += t.Next.ShowString()
	}
	return res
}

func (t *tType) ConvertibleTo(i interface{}) bool {
	u := MustToType(i)
	return t.Type.ConvertibleTo(u.RefType())
}

func (t *tType) Implements(i interface{}) bool {
	u := MustToType(i)
	return t.Type.Implements(u.RefType())
}

func (t *tType) FieldByName(name string) Field {
	if t.kind == reflect.Struct && t.fields != nil {
		for _, f := range t.fields {
			if f.GetName() == name {
				return f
			}
		}
	}
	{
		f, ok := t.Type.FieldByName(name)
		if !ok {
			return nil
		}
		return &tField{
			TNoteCode: TNoteCode{nil},
			Type:      NewType(f.Type),
			ReName:    f.Name,
			Tag:       string(f.Tag),
		}
	}
}

func (t *tType) FieldTypeByName(name string) (Type, bool) {
	if t.kind == reflect.Struct && t.fields != nil {
		f := t.FieldByName(name)
		if f != nil {
			return f.GetType(), true
		}
	}
	f, ok := t.Type.FieldByName(name)
	if !ok {
		return nil, false
	}
	return NewType(f.Type), true
}

func (t *tType) NumField() int {
	if t.kind == reflect.Struct && t.fields != nil {
		return len(t.GetFields())
	}
	return t.Type.NumField()
}

func (t *tType) Field(i int) Field {
	if t.kind == reflect.Struct && t.fields != nil {
		return t.GetFields()[i]
	}
	f := t.Type.Field(i)
	return NewField(f.Name, NewTypeI(f.Type), string(f.Tag))
}

func (t *tType) Kind() reflect.Kind {
	if t.kind != 0 {
		return t.kind
	}
	if strings.HasPrefix(t.Str, "[]") {
		return reflect.Slice
	}
	if strings.HasPrefix(t.Str, "*") {
		return reflect.Ptr
	}
	if t.Type != nil {
		return t.Type.Kind()
	}
	if t.kind == reflect.Struct && t.fields != nil {
		return reflect.Struct
	}
	if t.Type == nil && t.Str == "" && t.Named != "" && t.Next != nil {
		return t.Next.Kind()
	}
	return t.RefType().Kind()
}

func (t *tType) SetKind(v reflect.Kind) {
	t.kind = v
}

func (t *tType) MethodByName(name string) (reflect.Method, bool) {
	return t.Type.MethodByName(name)
}

func (t *tType) Name() string {
	if t.Type != nil && t.Type.Name() != "" {
		return t.Type.Name()
	}
	{
		// from named
		if t.GetNamed() != "" {
			return t.GetNamed()
		}
	}
	return t.String()
}

func (t *tType) GetNamed() string {
	return t.Named
}

func (t *tType) GetRowStr() string {
	return t.Str
}

func (t *tType) GetNext() Type {
	if t.Next != nil {
		return t.Next
	}
	if t.Type != nil && t.Type.Name() == "" && (t.Type.Kind() == reflect.Array ||
		t.Type.Kind() == reflect.Chan ||
		t.Type.Kind() == reflect.Map ||
		t.Type.Kind() == reflect.Pointer ||
		t.Type.Kind() == reflect.Slice) {
		return &tType{
			TNoteCode:   TNoteCode{nil},
			Type:        t.Type.Elem(),
			Pkg:         t.Pkg,
			Str:         "",
			Named:       "",
			Next:        nil,
			inReference: t.inReference,
			kind:        0,
			fields:      nil,
			funcs:       nil,
		}
	}
	return nil
}

func (t *tType) SetNamed(v string) {
	t.Named = v
}

func (t *tType) SetPkg(v string) {
	t.Pkg = v
}

func (t *tType) Zero() Value {
	return &tValue{
		TNoteCode:    TNoteCode{nil},
		IType:        t,
		Action:       ValueActionZero,
		Left:         nil,
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

func (t *tType) InterfaceForType() bool {
	return true
}

func (t *tType) InReference() bool {
	return t.inReference
}

func (t *tType) SetInReference(v bool) {
	t.inReference = v
}

func (t *tType) GetFields() []Field {
	return t.fields
}

func (t *tType) AddFields(fs []Field) {
	t.fields = append(t.fields, fs...)
}

func (t *tType) GetFuncs() []Func {
	return t.funcs
}

func (t *tType) FuncByName(name string) Func {
	for _, f := range t.funcs {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

func (t *tType) WarpNamed(named string) Type {
	return &tType{
		TNoteCode:   TNoteCode{nil},
		Type:        nil,
		Str:         "",
		Pkg:         "",
		Named:       named,
		Next:        t,
		inReference: t.inReference,
		kind:        t.kind,
		fields:      nil,
		funcs:       nil,
	}
}
