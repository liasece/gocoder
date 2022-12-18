package gocoder

// PtrChecker type
type PtrChecker interface {
	Codable

	GetCheckerValue() []Value
	GetHandlers() []Codable
	GetIfNotNil() bool

	AddChecker(vi ...interface{}) PtrChecker
	C(cs ...Codable) PtrChecker
	InterfaceForPtrChecker() bool
}

type tPtrChecker struct {
	CheckerValue []Value
	Handlers     []Codable
	IfNotNil     bool
}

func (t *tPtrChecker) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tPtrChecker) GetCheckerValue() []Value {
	return t.CheckerValue
}

func (t *tPtrChecker) GetHandlers() []Codable {
	return t.Handlers
}

func (t *tPtrChecker) GetIfNotNil() bool {
	return t.IfNotNil
}

func (t *tPtrChecker) AddChecker(vi ...interface{}) PtrChecker {
	vs := MustToValueList(vi...)
	t.CheckerValue = append(t.CheckerValue, vs...)
	return t
}

func (t *tPtrChecker) C(cs ...Codable) PtrChecker {
	t.Handlers = append(t.Handlers, cs...)
	return t
}

func (t *tPtrChecker) Else(cs ...Codable) PtrChecker {
	t.Handlers = append(t.Handlers, cs...)
	return t
}

func (t *tPtrChecker) InterfaceForPtrChecker() bool {
	return true
}
