package services

import (
	"gonga/app/Models"

	"gorm.io/gorm"
)

func LoadNestedComments(comment *Models.Comment, db *gorm.DB) {
	db.Preload("User").Preload("Likes").Where("parent_id = ?", comment.ID).Find(&comment.Childrens)
	for i := range comment.Childrens {
		LoadNestedComments(comment.Childrens[i], db)
	}
}
