package factory

import (
	"gonga/app/Models"
	"github.com/go-faker/faker/v4"
)

// PostFactory generates a fake post instance
func PostFactory() Models.Post {
	post := Models.Post{
		Title:        faker.Sentence(),
		Body:         faker.Paragraph(),
		UserID:       0, // Set the appropriate user ID here
		LikeCount:    0,
		CommentCount: 0,
		ViewCount:    0,
		ShareCount:   0,
		IsPromoted:   false,
		IsFeatured:   false,
		Visibility:   Models.VisibilityPublic,
		
	}

	return post
}
