package gocoder

type Interface interface {
	Codeable
	GetFuncs() []Func
	FuncByName(name string) Func
	GetName() string
	GetType() Type
	IsInterface()
}

var _ Interface = (*tInterface)(nil)

type tInterface struct {
	Funcs  []Func
	ReName string
}

func (t *tInterface) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tInterface) GetName() string {
	return t.ReName
}

func (t *tInterface) GetFuncs() []Func {
	return t.Funcs
}

func (t *tInterface) FuncByName(name string) Func {
	for _, f := range t.Funcs {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

func (t *tInterface) GetType() Type {
	return &tType{
		Str: t.GetName(),
	}
}

func (t *tInterface) IsInterface() {
}
