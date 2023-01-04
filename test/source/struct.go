package source_test

import (
	"context"
	"time"
)

/* BigStruct is a struct for test */
type BigStruct struct {
	ID      string `json:"id"` // this entity id
	SquadID string `bson:"squadID"`
	UserID  string `bson:"userID"`
	// this entity id
	CreateAt time.Time `bson:"createAt"`
	Ready    bool      `bson:"ready"`
}

type IBigStruct interface {
	GetBigStruct() BigStruct
	SetBigStruct(BigStruct) error
	/* SetBigStructByName set the name and type of the BigStruct */
	SetBigStructByName(ctx context.Context, name string, typ *BigStruct) (*BigStruct, bool, error)
}

// GetBigStruct func
func (b *BigStruct) GetBigStruct() BigStruct {
	return *b
}

/* SetBigStruct func */
func (b *BigStruct) IsBigStruct() bool {
	return true
}
