package gocoder

// ForRange func
type ForRange interface {
	Codable

	GetType() FuncType
	GetAutoSet() bool
	GetToValues() Value
	GetValue() Value
	GetCodes() []Codable

	C(cs ...Codable) ForRange
	ToCode() Code

	InterfaceForRange() bool
}

type tForRange struct {
	Type     FuncType
	AutoSet  bool
	ToValues Value
	Value    Value
	Codes    []Codable
}

func (t *tForRange) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tForRange) C(cs ...Codable) ForRange {
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

func (t *tForRange) GetCodes() []Codable {
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

func (t *tForRangeCode) C(codes ...Codable) Code {
	t.tForRange.C(codes...)
	return t
}
