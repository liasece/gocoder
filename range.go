package gocoder

// ForRange func
type ForRange interface {
	Codeable

	GetType() FuncType
	GetAutoSet() bool
	GetToValues() Value
	GetValue() Value
	GetCodes() []Codeable

	C(cs ...Codeable) ForRange
	ToCode() Code

	InterfaceForRange() bool
}

type tForRange struct {
	Type     FuncType
	AutoSet  bool
	ToValues Value
	Value    Value
	Codes    []Codeable
}

func (t *tForRange) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tForRange) C(cs ...Codeable) ForRange {
	t.Codes = append(t.Codes, cs...)
	return t
}

func (t *tForRange) GetType() FuncType {
	return t.Type
}

func (t *tForRange) GetAutoSet() bool {
	return t.AutoSet
}

func (t *tForRange) GetToValues() Value {
	return t.ToValues
}

func (t *tForRange) GetValue() Value {
	return t.Value
}

func (t *tForRange) GetCodes() []Codeable {
	return t.Codes
}

func (t *tForRange) ToCode() Code {
	return &tForRangeCode{
		tForRange: t,
	}
}

func (t *tForRange) InterfaceForRange() bool {
	return true
}

type tForRangeCode struct {
	*tForRange
}

func (t *tForRangeCode) C(codes ...Codeable) Code {
	t.tForRange.C(codes...)
	return t
}
