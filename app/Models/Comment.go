package Models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID   uint     `json:"-"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
	PostID   uint     `json:"-"`
	Post     Post     `json:"post" gorm:"foreignKey:PostID"`
	Body     string   `json:"body"`
	Likes    []Like   `json:"-" gorm:"foreignKey:CommentID"`
	ParentID *uint    `json:"-"`
	Parent   *Comment `json:"parent" gorm:"foreignKey:ParentID"`
	Mentions []*Mention `json:"mentions" gorm:"polymorphic:Owner;"`
}

func (Comment) TableName() string {
	return "comments"
}
