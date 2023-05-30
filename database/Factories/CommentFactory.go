package factory

import (
	"gonga/app/Models"

	faker "github.com/brianvoe/gofakeit/v6"
)

// CommentFactory generates a fake comment instance
func CommentFactory() Models.Comment {

	comment := Models.Comment{
		Body:     faker.Paragraph(2, 20, 20, "."),
		UserID:   0, // Set the appropriate user ID here
		PostID:   0, // Set the appropriate post ID here
		ParentID: nil,
	}

	return comment
}
