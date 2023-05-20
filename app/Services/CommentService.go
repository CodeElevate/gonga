package services

import (
	"gonga/app/Models"

	"gorm.io/gorm"
)

func LoadNestedComments(comment *Models.Comment, db *gorm.DB) {
	db.Preload("User").Preload("Likes").Where("parent_id = ?", comment.ID).Find(&comment.Children)
	for i := range comment.Children {
		LoadNestedComments(comment.Children[i], db)
	}
}
