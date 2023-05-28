package factory

import (
	"gonga/app/Models"
	"github.com/go-faker/faker/v4"
)

// CommentFactory generates a fake comment instance
func CommentFactory() Models.Comment {
	comment := Models.Comment{
		Body:     faker.Sentence(),
		UserID:   0, // Set the appropriate user ID here
		PostID:   0, // Set the appropriate post ID here
		ParentID: nil,
	}

	return comment
}
