package Models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID    uint       `json:"user_id"`
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PostID    uint       `json:"post_id"`
	Post      *Post      `json:"post,omitempty" gorm:"foreignKey:PostID;"`
	Body      string     `json:"body"`
	Likes     []Like     `json:"likes" gorm:"polymorphic:Likeable;"`
	ParentID  *uint      `json:"parent_id"`
	Parent    *Comment   `json:"parent,omitempty"`
	Childrens []*Comment `json:"childrens,omitempty" gorm:"foreignKey:ParentID"`
	Mentions  []*Mention `json:"mentions" gorm:"polymorphic:Owner;"`
}


func (Comment) TableName() string {
	return "comments"
}
