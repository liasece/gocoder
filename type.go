package gocoder

import (
	"reflect"
	"strings"
)

// Type type
type Type interface {
	Codable
	RefType() reflect.Type
	IsPtr() bool
	IsSlice() bool
	UnPtr() Type
	IsStruct() bool
	TackPtr() Type
	Slice() Type
	IsNil() bool
	Elem() Type
	Kind() reflect.Kind
	String() string
	ShowString() string
	CurrentCode() string
	Package() string
	PackageInReference() string
	ConvertibleTo(i interface{}) bool
	Implements(i interface{}) bool
	NumField() int
	Field(i int) Field
	FieldByName(name string) (reflect.StructField, bool)
	FieldTypeByName(name string) (Type, bool)
	MethodByName(name string) (reflect.Method, bool)
	Zero() Value
	Name() string
	GetNamed() string
	GetRowStr() string
	SetNamed(string)
	SetPkg(string)
	GetNext() Type
	GetStruct() Struct

	AllSub() []Type // list all type chian nodes, top type in last index
	HasPtrSubType() bool
	HasSliceSubType() bool
	InterfaceForType() bool
}

type tType struct {
	reflect.Type
	Str    string
	Pkg    string
	Struct Struct
	Named  string
	Next   Type
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
	// if t == nil || t.Type == nil {
	// 	return false
	// }
	if t.Type == nil && t.Str != "" {
		return t.Str[0] == '*'
	}
	return t.Kind() == reflect.Ptr
}

func (t *tType) IsSlice() bool {
	return t.Kind() == reflect.Slice
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
	InterfaceType = reflect.ValueOf(T{}).Field(0).Type()
}

func (t *tType) RefType() reflect.Type {
	if t.Type != nil {
		return t.Type
	}
	if t.Struct != nil {
		fields := make([]reflect.StructField, 0, len(t.Struct.GetFields()))
		for _, v := range t.Struct.GetFields() {
			fields = append(fields, reflect.StructField{
				Name: v.GetName(),
				Type: v.GetType().RefType(),
				Tag:  reflect.StructTag(v.GetTag()),
			})
		}
		return reflect.StructOf(fields)
	}
	return InterfaceType
}

func (t *tType) UnPtr() Type {
	// if t == nil || t.Type == nil {
	// 	return t
	// }
	if t.Kind() == reflect.Ptr {
		return &tType{
			Type:   t.Type.Elem(),
			Struct: t.Struct,
			Pkg:    t.Pkg,
		}
	}
	return t
}

func (t *tType) TackPtr() Type {
	if t.Type == nil {
		if !strings.HasPrefix(t.Str, "*") {
			return &tType{
				Str:    "*",
				Struct: t.Struct,
				Next:   t,
				// Pkg:    t.Pkg,
			}
		}
		return t
	}
	if t.Kind() != reflect.Ptr {
		return &tType{
			Type:   reflect.PtrTo(t.Type),
			Str:    t.Str,
			Struct: t.Struct,
			Pkg:    t.Pkg,
		}
	}
	return t
}

func (t *tType) Slice() Type {
	if t.Type == nil {
		return &tType{
			Str:    "[]",
			Struct: t.Struct,
			Next:   t,
			// Pkg:    t.Pkg,
		}
	}
	if t.Kind() != reflect.Ptr {
		str := t.Str
		if str != "" {
			str = "[]" + str
		}
		return &tType{
			Type:   reflect.SliceOf(t.Type),
			Struct: t.Struct,
			Str:    str,
		}
	}
	return t
}

func (t *tType) Elem() Type {
	if t.Str == "[]" || t.Str == "*" {
		tmp := *t
		res := &tmp
		res.Str = ""
		return res
	}
	tmp := *t
	res := &tmp
	res.Type = t.Type.Elem()
	return res
}

func (t *tType) Package() string {
	return t.Pkg
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
			res = res + t.Next.String()
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
		if t.Next != nil {
			res = t.Next.String()
		} else {
			res = t.Named
		}
	}
	return res
}

func (t *tType) ShowString() string {
	head := ""
	if t.Package() != "" {
		head = t.Package() + "."
	}
	if t.Str != "" {
		res := head + t.Str
		if t.Next != nil {
			res = res + t.Next.ShowString()
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
		res = res + t.Next.ShowString()
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

func (t *tType) FieldByName(name string) (reflect.StructField, bool) {
	return t.Type.FieldByName(name)
}

func (t *tType) FieldTypeByName(name string) (Type, bool) {
	// log.Error("test FieldByName", log.Any("name", name))
	if t.Struct != nil {
		f := t.Struct.FieldByName(name)
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
	if t.Struct != nil {
		return len(t.Struct.GetFields())
	}
	return t.Type.NumField()
}

func (t *tType) Field(i int) Field {
	if t.Struct != nil {
		return t.Struct.GetFields()[i]
	}
	f := t.Type.Field(i)
	return NewField(f.Name, NewTypeI(f.Type), string(f.Tag))
}

func (t *tType) Kind() reflect.Kind {
	if strings.HasPrefix(t.Str, "[]") {
		return reflect.Slice
	}
	if strings.HasPrefix(t.Str, "*") {
		return reflect.Ptr
	}
	if t.Type != nil {
		return t.Type.Kind()
	}
	if t.Struct != nil {
		return reflect.Struct
	}
	return t.RefType().Kind()
}

func (t *tType) MethodByName(name string) (reflect.Method, bool) {
	return t.Type.MethodByName(name)
}

func (t *tType) Name() string {
	if t.Type != nil {
		return t.Type.Name()
	}
	return t.GetNamed()
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
			Type:   t.Type.Elem(),
			Struct: t.Struct,
			Pkg:    t.Pkg,
		}
	}
	return nil
}

func (t *tType) SetNamed(v string) {
	t.Named = v
}

func (t *tType) SetPkg(v string) {
	// log.Debug("SetPkg", log.Any("v", v), log.Any("t.string", t.String()))
	t.Pkg = v
}

func (t *tType) Zero() Value {
	return &tValue{
		IType:  t,
		Action: ValueActionZero,
	}
}

func (t *tType) InterfaceForType() bool {
	return true
}

func (t *tType) GetStruct() Struct {
	return t.Struct
}
