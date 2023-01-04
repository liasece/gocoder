package ast

import (
	"go/ast"
	"strings"

	"github.com/liasece/gocoder"
)

func (c *CodeDecoder) GetNoteFromCommentGroup(_ DecoderContext, sts ...*ast.CommentGroup) []gocoder.Note {
	notes := make([]gocoder.Note, 0)
	for _, st := range sts {
		if st == nil {
			continue
		}
		for _, comment := range st.List {
			kind := gocoder.NoteKindLine
			text := strings.TrimSpace(comment.Text)
			if strings.HasPrefix(comment.Text, "//") {
				kind = gocoder.NoteKindLine
				text = text[2:]
			} else if strings.HasPrefix(comment.Text, "/*") {
				kind = gocoder.NoteKindBlock
				text = text[2 : len(text)-2]
			}
			text = strings.TrimSpace(text)
			notes = append(notes, gocoder.NewNote(text, kind))
		}
	}
	return notes
}
