package Models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID   uint       `json:"user_id"`
	User     User       `json:"user" gorm:"foreignKey:UserID"`
	PostID   uint       `json:"post_id"`
	Post     Post       `json:"post" gorm:"foreignKey:PostID"`
	Body     string     `json:"body"`
	Likes    []Like     `json:"likes" gorm:"foreignKey:CommentID"`
	ParentID *uint      `json:"parent_id"`
	Parent   *Comment   `json:"parent"`
	Children []*Comment `json:"childrens" gorm:"foreignKey:ID"`
	Mentions []*Mention `json:"mentions" gorm:"polymorphic:Owner;"`
}

func (Comment) TableName() string {
	return "comments"
}
