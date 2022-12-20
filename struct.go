package gocoder

type Struct interface {
	Codable
	GetFields() []Field
	AddFields([]Field)
	FieldByName(name string) Field
	GetName() string
	GetType() Type
	IsStruct()
}

var _ Struct = (*tStruct)(nil)

type tStruct struct {
	Fields []Field
	ReName string
}

func (t *tStruct) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tStruct) GetName() string {
	return t.ReName
}

func (t *tStruct) GetFields() []Field {
	return t.Fields
}

func (t *tStruct) AddFields(fs []Field) {
	t.Fields = append(t.Fields, fs...)
}

func (t *tStruct) FieldByName(name string) Field {
	for _, f := range t.Fields {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

func (t *tStruct) GetType() Type {
	return &tType{
		Str:    t.GetName(),
		Struct: t,
		Type:   nil,
		Pkg:    "",
		Named:  t.ReName,
		Next:   nil,
	}
}

func (t *tStruct) IsStruct() {
}
