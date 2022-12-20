package gocoder

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/liasece/log"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

// PkgTool interface
type PkgTool interface {
	PkgAlias(pkgPath string) string
	SetPkgAlias(pkgPath string, alias string)
	PkgAliasMap() map[string]string
}

// Writer type
type Writer interface {
	Line(is ...interface{})
	IsHead() bool
	Add(is ...interface{})
	AddCompact(is ...interface{})
	AddStr(strs ...interface{})
	AddNote(notes ...Note)
	Parentheses(is ...interface{})
	List(is ...interface{})
	ListValues(vs ...Value)
	ParenthesesValues(vs ...Value)
	ParenthesesArgs(vs ...Arg)
	ParenthesesTypes(vs ...Type)
	Block(is ...interface{})
	BlockCodes(cs ...Codable)
	InlineBlock(is ...interface{})
	InlineBlockCodes(cs ...Codable)
	WriteCode(c Codable)
	In()
	Out()
	InN(n int)
	OutN(n int)
	SetPkgTool(PkgTool)
}

// tWriter type
type tWriter struct {
	out        *bytes.Buffer
	indent     int
	notHead    bool
	needIndent bool
	inline     bool
	pkgTool    PkgTool
	toPkg      string
}

// ToCode func
func ToCode(c Codable, opts ...*ToCodeOption) string {
	opt := MergeToCodeOpt(opts...)
	pkgTool := opt.pkgTool
	if pkgTool == nil {
		pkgTool = NewDefaultPkgTool()
	}
	toPkg := ""
	if opt.pkgName != nil {
		toPkg = *opt.pkgName
	}
	if opt.pkgPath != nil {
		toPkg = *opt.pkgPath
	}
	w := &tWriter{
		out:        &bytes.Buffer{},
		pkgTool:    pkgTool,
		toPkg:      toPkg,
		indent:     0,
		notHead:    false,
		needIndent: false,
		inline:     false,
	}
	c.WriteCode(w)
	return w.out.String()
}

func GetImports(pkgTool PkgTool, skip []string) []string {
	m := pkgTool.PkgAliasMap()
	res := make([]string, 0)
	if len(m) > 0 {
		byAlias := make(map[string]string, len(m))
		aliases := make([]string, 0, len(m))

		for path, alias := range m {
			isSkip := false
			for _, s := range skip {
				if alias == s || path == s {
					isSkip = true
					break
				}
			}
			if isSkip {
				continue
			}
			aliases = append(aliases, alias)
			byAlias[alias] = path
		}

		sort.Strings(aliases)
		for _, alias := range aliases {
			if filepath.Base(byAlias[alias]) == alias {
				res = append(res, fmt.Sprintf("%q", byAlias[alias]))
			} else {
				res = append(res, fmt.Sprintf("%s %q", alias, byAlias[alias]))
			}
		}
	}
	return res
}

type ImportPkgStrSort []string

func (a ImportPkgStrSort) Len() int      { return len(a) }
func (a ImportPkgStrSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ImportPkgStrSort) Less(i, j int) bool {
	if strings.Count(a[i], ".") != 0 && strings.Count(a[j], ".") == 0 {
		return true
	}
	return strings.Count(a[i], "/") < strings.Count(a[j], "/")
}

func isInnerPkg(s string) bool {
	if strings.Count(s, ".") != 0 {
		return false
	}
	if strings.Count(s, "/") >= 3 {
		return false
	}
	return true
}

func isExternalPkg(s string) bool {
	return strings.Count(s, ".") != 0
}

func GetImportStr(pkgTool PkgTool, skip []string) string {
	imports := GetImports(pkgTool, skip)
	sort.Sort(ImportPkgStrSort(imports))
	if len(imports) > 0 {
		strs := ""
		for i, v := range imports {
			strs += v + "\n"
			if i < len(imports)-1 && ((isInnerPkg(imports[i]) && !isInnerPkg(imports[i+1])) || (isExternalPkg(imports[i]) && !isExternalPkg(imports[i+1]))) {
				strs += "\n"
			}
		}
		return fmt.Sprintf("import (\n%s\n)", strs)
	}
	return ""
}

func WriteToFileStr(c Codable, opts ...*ToCodeOption) (string, error) {
	buffer := &bytes.Buffer{}
	err := Write(buffer, c, opts...)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func Write(w io.Writer, c Codable, opts ...*ToCodeOption) error {
	opt := MergeToCodeOpt(opts...)
	if opt.pkgTool == nil {
		opt.pkgTool = NewDefaultPkgTool()
	}
	pkgName := "main"
	if opt.pkgName != nil {
		pkgName = *opt.pkgName
	}
	if pkgName != "" {
		_, err := fmt.Fprintln(w, "package", pkgName)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w)
		if err != nil {
			return err
		}
	}

	codeStr := ToCode(c, opt)

	importStr := GetImportStr(opt.pkgTool, []string{pkgName})
	if len(importStr) > 0 {
		_, err := fmt.Fprintf(w, "\n%s\n", importStr)
		if err != nil {
			return err
		}
	}
	{
		_, err := w.Write([]byte(codeStr))
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteToFile(filename string, c Codable, opts ...*ToCodeOption) error {
	opt := MergeToCodeOpt(opts...)
	if opt.pkgTool == nil {
		opt.pkgTool = NewDefaultPkgTool()
		if opt.pkgName != nil {
			opt.pkgTool.SetPkgAlias(*opt.pkgName, "")
		}
	}
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create directory")
	}
	str, err := WriteToFileStr(c, opt)
	if err != nil {
		return err
	}
	// log.L(nil).Error("test", log.Any("filename", filename), log.Any("str", str))
	bytes := []byte(str)
	if opt.noPretty == nil || !*opt.noPretty {
		bytes, err = imports.Process(filename, bytes, &imports.Options{
			FormatOnly: true,
			Comments:   true,
			TabIndent:  true,
			TabWidth:   8,
			Fragment:   false,
			AllErrors:  false,
		})
		if err != nil {
			err1 := os.WriteFile(filename, []byte(str), 0600)
			if err1 != nil {
				return errors.Wrapf(err1, "failed to write %s", filename)
			}
			return err
		}
	}
	err = os.WriteFile(filename, bytes, 0600)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s", filename)
	}

	return nil
}

// Line func
func (w *tWriter) SetPkgTool(v PkgTool) {
	w.pkgTool = v
}

// Line func
func (w *tWriter) WriteCode(c Codable) {
	switch t := c.(type) {
	case Receiver:
		w.Add("(", t.GetName(), " ", t.GetType(), ")")
	case Struct:
		w.Line("type ", t.GetName(), " struct {")
		fs := t.GetFields()
		is := make([]interface{}, len(fs))
		for i, v := range fs {
			is[i] = v
		}
		w.Add(is...)
		w.Line("}")
	case Interface:
		w.Line("type ", t.GetName(), " interface {")
		fs := t.GetFuncs()
		for _, v := range fs {
			w.InterfaceFuncToCode(v)
		}
		w.Line("}")
	case Field:
		is := []interface{}{t.GetName(), " ", t.GetType()}
		if t.GetTag() != "" {
			is = append(is, " `"+t.GetTag()+"`")
		}
		w.Line(is...)
	case Note:
		w.NoteToCode(t)
	case Type:
		str := typeStringOut(t, w.pkgTool, w.toPkg)
		// if str == "" {
		// 	log.Panic("typeStringOut str == \"\"", log.Reflect("t", t))
		// }
		if str == "" && t.GetNamed() != "" {
			w.AddStr(t.GetNamed() + " ")
		} else {
			w.Add(str, t.GetNext())
		}
	case Value:
		w.ValueToCode(t)
	case If:
		w.IfToCode(t)
	case ElseIf:
		w.IfToCode(t)
	case Else:
		w.IfToCode(t)
	case PtrChecker:
		w.PtrCheckerToCode(t)
	case Arg:
		if t.GetVariableLength() {
			w.Add(t.GetName(), " ...", t.GetType())
		} else {
			w.Add(t.GetName(), " ", t.GetType())
		}
	case Func:
		w.FuncToCode(t)
	case ForRange:
		w.ForRangeToCode(t)
	case Return:
		w.Add("return ", t.GetValue())
	case Code:
		w.CodeToCode(t)
	default:
		panic(fmt.Sprintf("WriteCode unknown type: %+v", c))
	}
}

func isValidPkgName(str string) bool {
	ok, err := regexp.Match("^[_a-zA-Z]([_a-zA-Z0-9])*$", []byte(str))
	if err != nil {
		log.Panic("isValidPkgName error", log.ErrorField(err))
	}
	return ok
}

// Line func
func (w *tWriter) ValueToCode(t Value) {
	if t.GetIType() != nil {
		typeStringOut(t.GetIType(), w.pkgTool, w.toPkg)
	}
	switch t.GetAction() {
	case ValueActionNone:
		switch {
		case t.GetName() != "":
			if t.GetLeft() != nil {
				w.Add(t.GetLeft())
				w.AddStr(".")
			}
			name := t.GetName()
			if li := strings.Split(name, "."); len(li) == 2 && isValidPkgName(li[0]) {
				ms := w.pkgTool.PkgAliasMap()
				find := false
				for _, v := range ms {
					if v == li[0] {
						find = true
						break
					}
				}
				if !find {
					pkg := w.pkgTool.PkgAlias(li[0])
					if pkg != "" {
						name = strings.Join([]string{pkg, li[1]}, ".")
					}
					log.Warn("value alias not find, use alias", log.Reflect("value", t))
				}
			}
			w.Add(name)
		case t.GetSrcValue() != nil:
			if str, ok := t.GetSrcValue().(string); ok {
				w.Add(`"` + str + `"`)
			} else {
				w.Add(t.GetSrcValue())
			}
		case t.GetValues() != nil:
			w.ListValues(t.GetValues()...)
		case t.Type() != nil:
			w.Add(t.Type())
		default:
			panic(fmt.Sprintf("unknown value: %+v", t))
		}
	case ValueActionZero:
		w.Add(getZeroValueCode(t.Type(), w.pkgTool))
	case ValueActionCastType:
		w.Parentheses(t.GetLeft())
		w.Parentheses(t.GetRight())
	case ValueActionAssertionType:
		w.Add(t.GetLeft())
		w.AddStr(".")
		w.Parentheses(t.Type())
	case ValueActionIndex:
		if lv, ok := t.GetLeft().(Value); ok && lv.NeedParent() {
			w.Parentheses(t.GetLeft())
		} else {
			w.Add(t.GetLeft())
		}
		w.AddStr("[")
		w.AddCompact(t.GetRight())
		w.AddStr("]")
	case ValueActionFuncCall:
		if f := t.GetFunc(); f != nil {
			w.WriteCode(f)
			w.ParenthesesValues(t.GetCallArgs()...)
		} else {
			w.Add(t.GetLeft())
			w.ParenthesesValues(t.GetCallArgs()...)
		}
	case ValueActionDot:
		w.Add(t.GetLeft())
		w.AddStr(".")
		w.Add(t.GetName())
	default:
		if t.GetLeft() != nil {
			if lv, ok := t.GetLeft().(Value); ok && lv.NeedParent() {
				w.Parentheses(t.GetLeft())
			} else {
				w.Add(t.GetLeft())
			}
		}
		if v, ok := actionCodeConv[t.GetAction()]; ok {
			w.Add(v)
		} else {
			if t.GetLeft() != nil {
				w.Add(" ")
			}
			w.Add(t.GetAction())
			if t.GetLeft() != nil {
				w.Add(" ")
			}
		}
		if t.GetRight() != nil {
			lv := t.GetRight()
			if lv.NeedParent() || (t.GetAction() == ValueActionNot && lv.Depth() > 1) {
				w.Parentheses(t.GetRight())
			} else {
				if t.GetLeft() == nil {
					w.AddCompact(t.GetRight())
				} else {
					w.Add(t.GetRight())
				}
			}
		}
	}
	for _, note := range t.GetNotes() {
		w.Add(note)
	}
}

func (w *tWriter) CodeToCode(t Code) {
	for index, v := range t.GetCodes() {
		if index != 0 {
			if _, ok := v.(Func); ok {
				w.Line()
			}
		}
		w.Add(v)
		if !w.IsHead() {
			w.Line()
		}
	}
}

func (w *tWriter) FuncToCode(t Func) {
	if t.GetType() == FuncTypeInline {
		oldInline := w.inline
		w.inline = true
		defer func() {
			w.inline = oldInline
		}()
	}
	if t.GetType() != FuncTypeInline {
		if len(t.GetNotes()) > 0 {
			w.AddNote(t.GetNotes()...)
		}
	}
	w.Add("func")
	if t.GetReceiver() != nil {
		w.Add(t.GetReceiver())
	}
	if t.GetName() != "" {
		w.AddStr(" ")
		w.Add(t.GetName())
	}
	w.ParenthesesArgs(t.GetArgs()...)
	if len(t.GetReturns()) > 1 {
		w.AddStr(" ")
		w.ParenthesesTypes(t.GetReturns()...)
	} else if len(t.GetReturns()) > 0 {
		w.AddStr(" ")
		w.Add(t.GetReturns()[0])
	}
	w.AddStr(" ")
	if t.GetType() == FuncTypeInline {
		w.InlineBlockCodes(t.GetCodes()...)
	} else {
		w.BlockCodes(t.GetCodes()...)
	}
	if t.GetType() == FuncTypeInline {
		if len(t.GetNotes()) > 0 {
			w.AddNote(t.GetNotes()...)
		}
	}
}

func (w *tWriter) InterfaceFuncToCode(t Func) {
	if t.GetType() == FuncTypeInline {
		panic("inline func can't be interface method")
	}
	if t.GetName() != "" {
		w.AddStr(" ")
		w.Add(t.GetName())
	}
	w.ParenthesesArgs(t.GetArgs()...)
	if len(t.GetReturns()) > 1 {
		w.AddStr(" ")
		w.ParenthesesTypes(t.GetReturns()...)
	} else if len(t.GetReturns()) > 0 {
		w.AddStr(" ")
		w.Add(t.GetReturns()[0])
	}
	w.Line()
}

func (w *tWriter) IfToCode(t BaseIf) {
	if t.Pre() != nil {
		w.Add("else ")
	}
	if t.GetValue() != nil {
		w.Add("if ", t.GetValue(), " ")
	}
	w.BlockCodes(t.GetCodes()...)

	if t.Next() != nil {
		w.AddStr(" ")
		t.Next().WriteCode(w)
	}
}

func (w *tWriter) NoteToCode(t Note) {
	if t.GetContent() != "" {
		switch t.GetKind() {
		case NoteKindNone:
		case NoteKindLine:
			if !w.IsHead() {
				w.AddStr(" ")
			}
			w.Line("// ", t.GetContent())
		case NoteKindBlock:
			w.Add(" /* ", t.GetContent(), " */")
		}
	}
}

func (w *tWriter) PtrCheckerToCode(t PtrChecker) {
	ptrValues := make([]Value, 0, len(t.GetCheckerValue()))
	for _, v := range t.GetCheckerValue() {
		if v.IsPtr() || v.Type() == nil || v.Type().Kind() == reflect.Slice || v.Type().Kind() == reflect.Map || v.Type().Kind() == reflect.Interface {
			ptrValues = append(ptrValues, v)
		}
	}
	if len(ptrValues) > 0 {
		var src Value
		if t.GetIfNotNil() {
			src = ptrValues[0].NE(NewValueNil())
		} else {
			src = ptrValues[0].Equal(NewValueNil())
		}
		for i := 1; i < len(ptrValues); i++ {
			if t.GetIfNotNil() {
				src = src.And(ptrValues[i].NE(NewValueNil()))
			} else {
				src = src.Or(ptrValues[i].Equal(NewValueNil()))
			}
		}
		code := NewIf(src, t.GetHandlers()...)
		w.Add(code)
	} else {
		log.Error("PtrCheckerToCode but target type not ptr type", log.Any("tValues", t.GetCheckerValue()))
		w.Add(NewCode().C(t.GetHandlers()...))
	}
}

func (w *tWriter) ForRangeToCode(t ForRange) {
	if t.GetType() == FuncTypeInline {
		oldInline := w.inline
		w.inline = true
		defer func() {
			w.inline = oldInline
		}()
	}
	w.Add("for ")
	w.Add(t.GetToValues())
	if t.GetAutoSet() {
		w.Add(" := ")
	} else {
		w.Add(" = ")
	}
	w.Add("range ")
	w.Add(t.GetValue().UnPtr())
	w.AddStr(" ")
	if t.GetType() == FuncTypeInline {
		w.InlineBlockCodes(t.GetCodes()...)
	} else {
		w.BlockCodes(t.GetCodes()...)
	}
}

// IsHead func
func (w *tWriter) IsHead() bool {
	return !w.notHead
}

// SetHead func
func (w *tWriter) SetHead(head bool) {
	w.notHead = !head
}

// Line func
func (w *tWriter) Line(is ...interface{}) {
	w.Add(is...)
	if !w.inline {
		w.AddStr("\n")
		w.SetHead(true)
		w.needIndent = true
	}
}

// Add func
func (w *tWriter) Add(is ...interface{}) {
	w.add(is...)
}

// AddCompact func
func (w *tWriter) AddCompact(is ...interface{}) {
	w.add(is...)
}

// IsNil check interface is nil, like IsNil((*int)(nil)) == true
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	defer func() {
		_ = recover()
	}()
	vi := reflect.ValueOf(i)
	return vi.IsNil()
}

// Add func
func (w *tWriter) add(is ...interface{}) {
	for _, i := range is {
		if IsNil(i) {
			continue
		}
		if w.IsHead() && w.needIndent {
			w.AddStr(strings.Repeat("\t", w.indent))
			w.needIndent = false
		}
		if v, ok := i.(Codable); ok {
			v.WriteCode(w)
		} else {
			w.AddStr(i)
		}
	}
}

// Parentheses func // (***)
func (w *tWriter) Parentheses(is ...interface{}) {
	w.list(true, is...)
}

// List func
func (w *tWriter) List(is ...interface{}) {
	w.list(false, is...)
}

// List func
func (w *tWriter) list(parent bool, is ...interface{}) {
	if len(is) == 0 && !parent {
		return
	}
	if w.IsHead() && w.needIndent {
		w.AddStr(strings.Repeat("\t", w.indent))
		w.needIndent = false
	}
	w.SetHead(false)
	if parent {
		w.AddStr("(")
	}
	prefix := ""
	for index, i := range is {
		if index != 0 {
			w.AddStr(", ")
		}
		if v, ok := i.(Codable); ok {
			w.WriteCode(v)
		} else {
			w.AddStr(prefix, i)
		}
		prefix = " "
	}
	if parent {
		w.AddStr(")")
	}
}

// Block func
func (w *tWriter) Block(is ...interface{}) {
	if !w.inline {
		w.Line("{")
		w.In()
	} else {
		w.Add("{ ")
	}
	for index, i := range is {
		if index != 0 && w.inline {
			w.AddStr("; ")
		}
		w.Add(i)
		if !w.IsHead() {
			w.Line()
		}
	}
	if !w.inline {
		w.Out()
		if !w.IsHead() {
			w.Line()
		}
		w.Add("}")
	} else {
		w.Add("}")
	}
}

// BlockCodes func
func (w *tWriter) BlockCodes(cs ...Codable) {
	is := make([]interface{}, len(cs))
	for i, v := range cs {
		is[i] = v
	}
	w.Block(is...)
}

// InlineBlock func
func (w *tWriter) InlineBlock(is ...interface{}) {
	oldInline := w.inline
	w.inline = true
	defer func() {
		w.inline = oldInline
	}()
	w.Block(is...)
}

// InlineBlockCodes func
func (w *tWriter) InlineBlockCodes(cs ...Codable) {
	is := make([]interface{}, len(cs))
	for i, v := range cs {
		is[i] = v
	}
	w.InlineBlock(is...)
}

// ListValues func
func (w *tWriter) ListValues(vs ...Value) {
	is := make([]interface{}, len(vs))
	for i, v := range vs {
		is[i] = v
	}
	w.List(is...)
}

// ParenthesesValues func
func (w *tWriter) ParenthesesValues(vs ...Value) {
	is := make([]interface{}, len(vs))
	for i, v := range vs {
		is[i] = v
	}
	w.Parentheses(is...)
}

// ParenthesesArgs func
func (w *tWriter) ParenthesesArgs(vs ...Arg) {
	is := make([]interface{}, len(vs))
	for i, v := range vs {
		is[i] = v
	}
	w.Parentheses(is...)
}

// ParenthesesTypes func
func (w *tWriter) ParenthesesTypes(vs ...Type) {
	is := make([]interface{}, len(vs))
	for i, v := range vs {
		is[i] = v
	}
	w.Parentheses(is...)
}

// AddStr func
func (w *tWriter) AddStr(strs ...interface{}) {
	noNil := make([]interface{}, 0, len(strs))
	for _, v := range strs {
		if v != nil {
			noNil = append(noNil, v)
		}
	}
	strs = noNil
	str := fmt.Sprint(strs...)
	if str != "" {
		if w.IsHead() {
			tailLineStr := str
			if index := strings.Index(str, "\n"); index >= 0 {
				tailLineStr = str[index:]
			}
			head := true
			for _, r := range tailLineStr {
				if r == ' ' || r == '\t' || r == '\n' {
					continue
				}
				head = false
				break
			}
			w.SetHead(head)
		}
	}
	fmt.Fprint(w.out, strs...)
}

// AddNote func
func (w *tWriter) AddNote(notes ...Note) {
	is := make([]interface{}, len(notes))
	for i, v := range notes {
		is[i] = v
	}
	w.Add(is...)
}

// In func
func (w *tWriter) In() {
	w.indent++
}

// Out func
func (w *tWriter) Out() {
	if w.indent-1 < 0 {
		log.Fatal("w.indent <= 0")
	}
	w.indent--
}

// InN func
func (w *tWriter) InN(n int) {
	w.indent += n
}

// OutN func
func (w *tWriter) OutN(n int) {
	if w.indent-n < 0 {
		log.Fatal("w.indent-n <= 0")
	}
	w.indent -= n
}
