package gocoder

// ElseIf type
type ElseIf interface {
	BaseIf
	Else(...Codeable) Else
	C(cs ...Codeable) ElseIf
}

// Else type
type Else interface {
	BaseIf
	C(cs ...Codeable) Else
}

// If type
type If interface {
	BaseIf
	Else(...Codeable) Else
	ElseIf(interface{}, ...Codeable) ElseIf
	C(cs ...Codeable) If
	ToCode() Code
}

// BaseIf type
type BaseIf interface {
	Codeable
	SetPre(pre BaseIf)
	SetNext(next BaseIf)
	Append(base BaseIf)
	Tail() BaseIf
	Next() BaseIf
	Pre() BaseIf
	AddCode(cs ...Codeable)
	GetValue() Value
	GetCodes() []Codeable
	InterfaceForIf() bool
}

type tIf struct {
	IfV   Value
	Codes []Codeable
	IPre  BaseIf
	INext BaseIf
}

func (t *tIf) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tIf) Else(codes ...Codeable) Else {
	t.Tail().Append(&tIf{
		Codes: codes,
	})
	return &tElse{
		tIf: t,
	}
}

func (t *tIf) C(codes ...Codeable) If {
	t.Tail().AddCode(codes...)
	return t
}

func (t *tIf) AddCode(codes ...Codeable) {
	t.Codes = append(t.Codes, codes...)
}

func (t *tIf) ElseIf(i interface{}, codes ...Codeable) ElseIf {
	v := MustToValue("", i)
	t.Tail().Append(&tIf{
		IfV:   v,
		Codes: codes,
	})
	return &tElseIf{
		tIf: t,
	}
}

func (t *tIf) Append(base BaseIf) {
	p := t.Tail()
	p.SetNext(base)
	base.SetPre(p)
}

func (t *tIf) Tail() BaseIf {
	var p BaseIf = t
	for p.Next() != nil {
		p = p.Next()
	}
	return p
}

func (t *tIf) Pre() BaseIf {
	return t.IPre
}

func (t *tIf) Next() BaseIf {
	return t.INext
}

func (t *tIf) SetNext(next BaseIf) {
	t.INext = next
}

func (t *tIf) SetPre(pre BaseIf) {
	t.IPre = pre
}

func (t *tIf) GetCodes() []Codeable {
	return t.Codes
}

func (t *tIf) GetValue() Value {
	return t.IfV
}

func (t *tIf) ToCode() Code {
	return &tIfCode{
		tIf: t,
	}
}

func (t *tIf) InterfaceForIf() bool {
	return true
}

type tElseIf struct {
	*tIf
}

func (t *tElseIf) C(codes ...Codeable) ElseIf {
	t.Tail().AddCode(codes...)
	return t
}

type tElse struct {
	*tIf
}

func (t *tElse) C(codes ...Codeable) Else {
	t.Tail().AddCode(codes...)
	return t
}

type tIfCode struct {
	*tIf
}

func (t *tIfCode) C(codes ...Codeable) Code {
	t.Tail().AddCode(codes...)
	return t
}
