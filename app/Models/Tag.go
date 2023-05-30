package Models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Title        string  `json:"title" gorm:"unique" faker:"unique"`
	CoverImage   string  `json:"cover_image"`
	BackendImage string  `json:"backend_image"`
	Description  string  `json:"description"`
	Color        string  `json:"color"`
	Slug         string  `json:"slug"`
	UserID       uint    `json:"user_id"`
	User         User    `json:"user" gorm:"foreignKey:UserID"`
	// ParentID     uint    `json:"-"`
	// Parent       *Tag    `json:"-" gorm:"foreignKey:ParentID"`
	// Children     []*Tag  `json:"-" gorm:"foreignKey:ParentID"`
	Posts        []*Post `json:"posts" gorm:"many2many:post_hashtags;"`
	// Interests []*Interest `json:"interests" gorm:"many2many:interest_tags;"`
}

func (Tag) TableName() string {
	return "tags"
}
