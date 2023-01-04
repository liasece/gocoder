package gocoder

type Field interface {
	Codable
	NoteCode

	GetTag() string
	GetName() string
	GetType() Type
	IsField()

	Clone() Field
}

var _ Field = (*tField)(nil)

type tField struct {
	TNoteCode
	Type   Type
	ReName string
	Tag    string
}

func (t *tField) Clone() Field {
	res := &tField{
		TNoteCode: t.TNoteCode.Clone(),
		Type:      t.Type,
		ReName:    t.ReName,
		Tag:       t.Tag,
	}
	if t.Type != nil {
		res.Type = t.Type.Clone()
	}
	return res
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
