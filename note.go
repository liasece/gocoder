package gocoder

// Note type
type Note interface {
	Codable
	GetContent() string
	GetKind() NoteKind
	InterfaceForNote() bool
	Clone() Note
}

type NoteCode interface {
	Notes() []Note
	SetNotes([]Note)
	AddNotes(...Note)
}

var _ NoteCode = (*TNoteCode)(nil)

type TNoteCode struct {
	notes []Note
}

func (t *TNoteCode) Clone() TNoteCode {
	res := TNoteCode{
		notes: nil,
	}
	if t.notes != nil {
		res.notes = make([]Note, len(t.notes))
		for i, n := range t.notes {
			res.notes[i] = n.Clone()
		}
	}
	return res
}

func (t *TNoteCode) Notes() []Note {
	return t.notes
}

func (t *TNoteCode) SetNotes(notes []Note) {
	t.notes = notes
}

func (t *TNoteCode) AddNotes(notes ...Note) {
	t.notes = append(t.notes, notes...)
}

// NoteKind type
type NoteKind int

// NoteKind type
const (
	NoteKindNone  NoteKind = 0
	NoteKindLine  NoteKind = 1
	NoteKindBlock NoteKind = 2
)

var _ Note = (*tNote)(nil)

type tNote struct {
	Content string
	Kind    NoteKind
}

func (t *tNote) Clone() Note {
	return &tNote{
		Content: t.Content,
		Kind:    t.Kind,
	}
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
