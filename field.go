package gocoder

type Field interface {
	Codeable
	Type
	GetTag() string
	GetName() string
	GetType() Type
	IsField()
}

var _ Field = (*tField)(nil)

type tField struct {
	Type
	ReName string
	Tag    string
}

func (t *tField) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tField) GetTag() string {
	return t.Tag
}

func (t *tField) GetName() string {
	return t.ReName
}

func (t *tField) GetType() Type {
	return t.Type
}

func (t *tField) IsField() {
}
