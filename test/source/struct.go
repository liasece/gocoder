package source

import "time"

type BigStruct struct {
	ID       string    `json:"id"`
	SquadID  string    `bson:"squadID"`
	UserID   string    `bson:"userID"`
	CreateAt time.Time `bson:"createAt"`
	Ready    bool      `bson:"ready"`
}
