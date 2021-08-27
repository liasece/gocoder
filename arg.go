package gocoder

// Arg type
type Arg interface {
	Codeable

	// Getter
	GetName() string
	GetType() Type
	GetValue() Value

	InterfaceForArg() bool
}

type tArg struct {
	Name string
	Type Type
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

func (t *tArg) GetValue() Value {
	return NewValue(t.Name, t.Type)
}
