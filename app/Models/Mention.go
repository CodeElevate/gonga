package Models

import (
	"gorm.io/gorm"
)

type Mention struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	User      User   `json:"user"`
	OwnerID   uint   `json:"owner_id"`
	OwnerType string `json:"owner_type"` // posts, comments, etc.
	Position  int    `json:"position"`   // position of the mention in the content
}

func (Mention) TableName() string {
	return "mentions"
}
