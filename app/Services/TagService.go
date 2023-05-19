package services

import (
	"errors"
	"gonga/app/Models"

	"gorm.io/gorm"
)

func EditTags(db *gorm.DB, postID string, hashtags []Models.Tag, userID uint) error {
	var post Models.Post
	result := db.Preload("Hashtags").First(&post, postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		} else {
			return result.Error
		}
	}
	if post.UserID != userID {
		return errors.New("you are not authorized to update post")
	}
	existingTags := make(map[string]*Models.Tag)
	for _, tag := range post.Hashtags {
		existingTags[tag.Title] = tag
	}

	updatedTags := make([]*Models.Tag, 0)
	removedTags := make([]*Models.Tag, 0)

	for _, hashtag := range hashtags {
		if tag, ok := existingTags[hashtag.Title]; ok {
			updatedTags = append(updatedTags, tag)
			delete(existingTags, hashtag.Title)
		} else {
			// Check if the tag already exists in the tags table
			var existingTag Models.Tag
			result := db.First(&existingTag, "title = ?", hashtag.Title)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					// Tag doesn't exist, create a new record for it
					newTag := &Models.Tag{Title: hashtag.Title, UserID: userID}
					result := db.Create(newTag)
					if result.Error != nil {
						return result.Error
					}
					updatedTags = append(updatedTags, newTag)
				} else {
					return result.Error
				}
			} else {
				// Tag already exists, use the existing tag record
				updatedTags = append(updatedTags, &existingTag)
			}
		}
	}

	for _, tag := range existingTags {
		removedTags = append(removedTags, tag)
	}

	if len(updatedTags) > 0 {
		err := db.Model(&post).Association("Hashtags").Replace(updatedTags)
		if err != nil {
			return err
		}
	}

	if len(removedTags) > 0 {
		err := db.Model(&post).Association("Hashtags").Delete(removedTags)
		if err != nil {
			return err
		}
	}

	return nil
}
