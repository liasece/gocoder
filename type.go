package gocoder

import (
	"reflect"
	"strings"
)

// Type type
type Type interface {
	Codeable
	RefType() reflect.Type
	IsPtr() bool
	UnPtr() Type
	IsStruct() bool
	TackPtr() Type
	Slice() Type
	IsNil() bool
	Elem() Type
	Kind() reflect.Kind
	String() string
	Package() string
	ConvertibleTo(i interface{}) bool
	Implements(i interface{}) bool
	FieldByName(name string) (reflect.StructField, bool)
	FieldTypeByName(name string) (Type, bool)
	MethodByName(name string) (reflect.Method, bool)
	Zero() Value
	Name() string
	GetNamed() string
	SetNamed(string)
	GetNext() Type

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
	return t.Kind() == reflect.Ptr
}

func (t *tType) RefType() reflect.Type {
	return t.Type
}

func (t *tType) UnPtr() Type {
	// if t == nil || t.Type == nil {
	// 	return t
	// }
	if t.Kind() == reflect.Ptr {
		return &tType{
			Type:   t.Type.Elem(),
			Struct: t.Struct,
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
			}
		}
		return t
	}
	if t.Kind() != reflect.Ptr {
		return &tType{
			Type:   reflect.PtrTo(t.Type),
			Struct: t.Struct,
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
		}
	}
	if t.Kind() != reflect.Ptr {
		return &tType{
			Type:   reflect.SliceOf(t.Type),
			Struct: t.Struct,
		}
	}
	return t
}

func (t *tType) Elem() Type {
	tmp := *t
	res := &tmp
	res.Type = t.Type.Elem()
	return res
}

func (t *tType) Package() string {
	return t.Pkg
}

func (t *tType) String() string {
	res := ""
	if t.Str != "" {
		res = t.Str
	}
	if t.Type != nil {
		str := t.Type.String()
		if str == "[]uint8" {
			str = "[]byte"
		}
		res += str
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
	if t.Struct != nil {
		f := t.Struct.FieldByName(name)
		if f != nil {
			return f, true
		}
	}
	f, ok := t.Type.FieldByName(name)
	if !ok {
		return nil, false
	}
	return NewType(f.Type), true
}

func (t *tType) MethodByName(name string) (reflect.Method, bool) {
	return t.Type.MethodByName(name)
}

func (t *tType) Name() string {
	return t.Type.Name()
}

func (t *tType) GetNamed() string {
	return t.Named
}

func (t *tType) GetNext() Type {
	return t.Next
}

func (t *tType) SetNamed(v string) {
	t.Named = v
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
