package source

import (
	"context"
	"time"
)

type BigStruct struct {
	ID       string    `json:"id"`
	SquadID  string    `bson:"squadID"`
	UserID   string    `bson:"userID"`
	CreateAt time.Time `bson:"createAt"`
	Ready    bool      `bson:"ready"`
}

type IBigStruct interface {
	GetBigStruct() BigStruct
	SetBigStruct(BigStruct) error
	SetBigStructByName(ctx context.Context, name string, typ *BigStruct) (*BigStruct, bool, error)
}

func (b *BigStruct) GetBigStruct() BigStruct {
	return *b
}

func (b *BigStruct) IsBigStruct() bool {
	return true
}
