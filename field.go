package gocoder

type Field interface {
	Codeable
	GetTag() string
	GetName() string
	GetType() Type
	IsField()
}

var _ Field = (*tField)(nil)

type tField struct {
	Type   Type
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
