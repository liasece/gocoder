package gocoder

// Return type
type Return interface {
	Codeable

	GetValue() Value
	InterfaceForReturn() bool
}

type tReturn struct {
	Value Value
}

func (t *tReturn) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tReturn) GetValue() Value {
	return t.Value
}

func (t *tReturn) InterfaceForReturn() bool {
	return true
}
