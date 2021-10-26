package gocoder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/liasece/log"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

// PkgTool interface
type PkgTool interface {
	PkgAlias(pkgPath string) string
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
	BlockCodes(cs ...Codeable)
	InlineBlock(is ...interface{})
	InlineBlockCodes(cs ...Codeable)
	WriteCode(c Codeable)
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
}

// ToCode func
func ToCode(c Codeable, opts ...*ToCodeOption) string {
	opt := MergeToCodeOpt(opts...)
	pkgTool := opt.pkgTool
	if pkgTool == nil {
		pkgTool = NewDefaultPkgTool()
	}
	w := &tWriter{
		out:     &bytes.Buffer{},
		pkgTool: pkgTool,
	}
	c.WriteCode(w)
	return w.out.String()
}

func WriteToFileStr(filename string, c Codeable, opts ...*ToCodeOption) (string, error) {
	opt := MergeToCodeOpt(opts...)
	if opt.pkgTool == nil {
		opt.pkgTool = NewDefaultPkgTool()
	}
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return "", errors.Wrap(err, "failed to create directory")
	}
	buffer := &bytes.Buffer{}
	pkgName := "main"
	if opt.pkgName != nil {
		pkgName = *opt.pkgName
	}
	_, err = fmt.Fprintln(buffer, "package", pkgName)
	if err != nil {
		return "", err
	}

	_, err = fmt.Fprintln(buffer)
	if err != nil {
		return "", err
	}

	codeStr := ToCode(c, opt)

	if len(opt.pkgTool.PkgAliasMap()) > 0 {
		byAlias := make(map[string]string, len(opt.pkgTool.PkgAliasMap()))
		aliases := make([]string, 0, len(opt.pkgTool.PkgAliasMap()))

		for path, alias := range opt.pkgTool.PkgAliasMap() {
			aliases = append(aliases, alias)
			byAlias[alias] = path
		}

		sort.Strings(aliases)
		_, err = fmt.Fprintln(buffer, "import (")
		if err != nil {
			return "", err
		}
		for _, alias := range aliases {
			_, err = fmt.Fprintf(buffer, "\t%s %q\n", alias, byAlias[alias])
			if err != nil {
				return "", err
			}
		}

		_, err = fmt.Fprintln(buffer, ")")
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(buffer, "")
		if err != nil {
			return "", err
		}
	}

	_, err = buffer.WriteString(codeStr)
	if err != nil {
		return "", err
	}

	bytes := buffer.Bytes()
	if opt.noPretty == nil || *opt.noPretty == false {
		bytes, err = imports.Process(filename, bytes, &imports.Options{FormatOnly: true, Comments: true, TabIndent: true, TabWidth: 8})
		if err != nil {
			err1 := ioutil.WriteFile(filename, buffer.Bytes(), 0644)
			if err1 != nil {
				return "", errors.Wrapf(err1, "failed to write %s", filename)
			}
			return "", err
		}
	}
	return string(bytes), nil
}

func WriteToFile(filename string, c Codeable, opts ...*ToCodeOption) error {
	str, err := WriteToFileStr(filename, c, opts...)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, []byte(str), 0644)
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
func (w *tWriter) WriteCode(c Codeable) {
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
	case Field:
		is := []interface{}{t.GetName(), " ", t.GetType()}
		if t.GetTag() != "" {
			is = append(is, " `"+t.GetTag()+"`")
		}
		w.Line(is...)
	case Note:
		w.NoteToCode(t)
	case Type:
		str := typeStringOut(t, w.pkgTool)
		if str == "" {
			panic("typeStringOut str == \"\"")
		}
		if t.GetNamed() != "" {
			w.AddStr(t.GetNamed() + " ")
		}
		w.Add(str)
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
		w.Add(t.GetName(), " ", t.GetType())
	case Func:
		w.FuncToCode(t)
	case ForRange:
		w.ForRangeToCode(t)
	case Return:
		w.Add("return ", t.GetValue())
	case Code:
		w.CodeToCode(t)
	default:
		panic(fmt.Errorf("WriteCode unknown type: %+v", c))
	}
}

// Line func
func (w *tWriter) ValueToCode(t Value) {
	if t.GetIType() != nil {
		typeStringOut(t.GetIType(), w.pkgTool)
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
			if li := strings.Split(name, "."); len(li) == 2 {
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
			w.Add(typeStringOut(t.Type(), w.pkgTool))
		default:
			panic(fmt.Errorf("unknown value: %+v", t))
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
		if len(t.GetNotes()) == 0 || !strings.HasPrefix(t.GetNotes()[0].GetContent(), t.GetName()) {
			w.Add(NewNote(t.GetName()+" func", NoteKindLine))
		}
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
		if v.IsPtr() {
			ptrValues = append(ptrValues, v)
		}
	}
	if len(ptrValues) > 0 {
		src := ptrValues[0].Equal(NewValueNil())
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
	w.add(true, is...)
}

// AddCompact func
func (w *tWriter) AddCompact(is ...interface{}) {
	w.add(false, is...)
}

// Add func
func (w *tWriter) add(interval bool, is ...interface{}) {
	for _, i := range is {
		if w.IsHead() && w.needIndent {
			w.AddStr(strings.Repeat("\t", w.indent))
			w.needIndent = false
		}
		if v, ok := i.(Codeable); ok {
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
		if v, ok := i.(Codeable); ok {
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
func (w *tWriter) BlockCodes(cs ...Codeable) {
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
func (w *tWriter) InlineBlockCodes(cs ...Codeable) {
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