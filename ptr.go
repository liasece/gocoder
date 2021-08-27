package gocoder

// PtrChecker type
type PtrChecker interface {
	Codeable

	GetCheckerValue() []Value
	GetHandlers() []Codeable
	GetIfNotNil() bool

	AddChecker(vi ...interface{}) PtrChecker
	C(cs ...Codeable) PtrChecker
	InterfaceForPtrChecker() bool
}

type tPtrChecker struct {
	CheckerValue []Value
	Handlers     []Codeable
	IfNotNil     bool
}

func (t *tPtrChecker) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tPtrChecker) GetCheckerValue() []Value {
	return t.CheckerValue
}

func (t *tPtrChecker) GetHandlers() []Codeable {
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

func (t *tPtrChecker) C(cs ...Codeable) PtrChecker {
	t.Handlers = append(t.Handlers, cs...)
	return t
}

func (t *tPtrChecker) Else(cs ...Codeable) PtrChecker {
	t.Handlers = append(t.Handlers, cs...)
	return t
}

func (t *tPtrChecker) InterfaceForPtrChecker() bool {
	return true
}
