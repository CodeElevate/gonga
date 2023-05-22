package Models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID       uint   `json:"user_id"`
	User         *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	LikeableID   uint   `json:"likable_id"`
	LikeableType string `json:"likable_type"` // posts, comments, users, etc.
}

func (Like) TableName() string {
	return "likes"
}
