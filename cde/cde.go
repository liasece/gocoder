package cde

import (
	"fmt"
	"reflect"

	"github.com/liasece/gocoder"
)

// Make func
func Make(typI interface{}, argsI ...interface{}) gocoder.Value {
	typ := Type(typI)
	args := gocoder.MustToValueList(argsI...)
	if typ.RefType() != nil {
		typ = typ.UnPtr()
		if typ.Kind() == reflect.Map && len(args) == 2 {
			panic(fmt.Errorf("make map func args must < 2"))
		}
	}
	if len(args) > 2 {
		panic(fmt.Errorf("make func args must < 2"))
	}
	args = append([]gocoder.Value{gocoder.NewOnlyTypeValue(typ)}, args...)
	return gocoder.NewValue("make", typ).Call(args)
}

// Len func
func Len(arg gocoder.Value) gocoder.Value {
	return gocoder.NewValueNameI("len", int(0)).Call(arg.UnPtr())
}

// ForRange func
func ForRange(autoSet bool, toValues gocoder.Value, value gocoder.Value, cs ...gocoder.Codeable) gocoder.ForRange {
	return gocoder.NewForRange(autoSet, gocoder.FuncTypeDefault, toValues, value)
}

// NoteLine func
func NoteLine(content string) gocoder.Note {
	return gocoder.MustToNote(gocoder.NoteKindLine, content)
}

// Note func
func Note(content string) gocoder.Note {
	return gocoder.MustToNote(gocoder.NoteKindBlock, content)
}

// PtrCheckerNil func
func PtrCheckerNil(checkerValue ...gocoder.Value) gocoder.PtrChecker {
	return gocoder.NewPtrChecker(false, checkerValue...)
}

// PtrCheckerNotNil func
func PtrCheckerNotNil(checkerValue ...gocoder.Value) gocoder.PtrChecker {
	return gocoder.NewPtrChecker(true, checkerValue...)
}

// If func
func If(v gocoder.Value, cs ...gocoder.Codeable) gocoder.If {
	return gocoder.NewIf(v, cs...)
}

// Type func
func Type(i interface{}) gocoder.Type {
	return gocoder.MustToType(i)
}

// TypeD func
func TypeD(pkg string, str string) gocoder.Type {
	return gocoder.NewTypeDetail(pkg, str)
}

// TypeError func
func TypeError() gocoder.Type {
	return gocoder.NewTypeName("error")
}

// Values func
func Values(is ...interface{}) gocoder.Value {
	return gocoder.MustToValues(is...)
}

// Value func
func Value(name string, i interface{}) gocoder.Value {
	return gocoder.MustToValue(name, i)
}

// Return func
func Return(is ...interface{}) gocoder.Return {
	if len(is) == 0 {
		return gocoder.NewReturn(nil)
	}
	if len(is) == 1 {
		return gocoder.NewReturn(Value("", is[0]))
	}
	return gocoder.NewReturn(Values(is...))
}

// Receiver func
func Receiver(name string, typ gocoder.Type) gocoder.Receiver {
	return gocoder.NewReceiver(name, typ)
}

// Method func
func Method(name string, receiver gocoder.Receiver, args []gocoder.Arg, returns []gocoder.Type, notes ...gocoder.Note) gocoder.Func {
	return gocoder.NewFunc(gocoder.FuncTypeDefault, name, receiver, args, returns, notes...)
}

// Func func
func Func(name string, args []gocoder.Arg, returns []gocoder.Type, notes ...gocoder.Note) gocoder.Func {
	return gocoder.NewFunc(gocoder.FuncTypeDefault, name, nil, args, returns, notes...)
}

// FuncInline func
func FuncInline(name string, args []gocoder.Arg, returns []gocoder.Type, notes ...gocoder.Note) gocoder.Func {
	return gocoder.NewFunc(gocoder.FuncTypeInline, name, nil, args, returns, notes...)
}

// Struct func
func Struct(name string, fs ...gocoder.Field) gocoder.Struct {
	return gocoder.NewStruct(name, fs)
}

// Interface func
func Interface(name string, fs ...gocoder.Func) gocoder.Interface {
	return gocoder.NewInterface(name, fs)
}

// Field func
func Field(name string, typ interface{}, tag string) gocoder.Field {
	return gocoder.NewField(name, Type(typ), tag)
}

// Arg func
func Arg(name string, i interface{}) gocoder.Arg {
	return gocoder.NewArg(name, Type(i), false)
}

// ArgVar func
func ArgVar(name string, i interface{}) gocoder.Arg {
	return gocoder.NewArg(name, Type(i), true)
}

// Args func
func Args(argsI ...interface{}) []gocoder.Arg {
	args := make([]gocoder.Arg, 0, len(argsI))
	for _, v := range argsI {
		switch v := v.(type) {
		case []gocoder.Arg:
			args = append(args, v...)
		case gocoder.Arg:
			args = append(args, v)
		case gocoder.Value:
			args = append(args, v.ToArg())
		case []gocoder.Value:
			for _, v := range v {
				args = append(args, v.ToArg())
			}
		case gocoder.Type:
			args = append(args, gocoder.NewArg("", v, false))
		case []gocoder.Type:
			for _, v := range v {
				args = append(args, gocoder.NewArg("", v, false))
			}
		default:
			panic(fmt.Sprint("Args unknown type: ", reflect.TypeOf(v).String()))
		}
	}
	return args
}

// Returns func
func Returns(returns ...gocoder.Type) []gocoder.Type {
	return returns
}
