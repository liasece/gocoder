package gocoder

// Note type
type Note interface {
	Codable
	GetContent() string
	GetKind() NoteKind
	InterfaceForNote() bool
}

// NoteKind type
type NoteKind int

// NoteKind type
const (
	NoteKindNone  NoteKind = 0
	NoteKindLine  NoteKind = 1
	NoteKindBlock NoteKind = 2
)

type tNote struct {
	Content string
	Kind    NoteKind
}

func (t *tNote) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tNote) GetContent() string {
	return t.Content
}

func (t *tNote) GetKind() NoteKind {
	return t.Kind
}

func (t *tNote) InterfaceForNote() bool {
	return true
}
