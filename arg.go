package gocoder

// Arg type
type Arg interface {
	Codable
	NoteCode

	// Getter
	GetName() string
	GetType() Type
	GetValue() Value
	GetVariableLength() bool

	InterfaceForArg() bool
}

var _ Arg = (*tArg)(nil)

type tArg struct {
	TNoteCode

	Name           string
	Type           Type
	VariableLength bool
}

func (t *tArg) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tArg) InterfaceForArg() bool {
	return false
}

func (t *tArg) GetName() string {
	return t.Name
}

func (t *tArg) GetType() Type {
	return t.Type
}

func (t *tArg) GetVariableLength() bool {
	return t.VariableLength
}

func (t *tArg) GetValue() Value {
	return NewValue(t.Name, t.Type)
}
