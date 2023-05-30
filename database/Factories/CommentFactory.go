package factory

import (
	"gonga/app/Models"

	"github.com/jaswdr/faker"
)

// CommentFactory generates a fake comment instance
func CommentFactory() Models.Comment {
	faker := faker.New()

	comment := Models.Comment{
		Body:     faker.Lorem().Paragraph(10),
		UserID:   0, // Set the appropriate user ID here
		PostID:   0, // Set the appropriate post ID here
		ParentID: nil,
	}

	return comment
}
