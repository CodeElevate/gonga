package factory

import (
	"gonga/app/Models"
	"time"

	faker "github.com/brianvoe/gofakeit/v6"
)

// PostFactory generates a fake post instance
func PostFactory() Models.Post {

	post := Models.Post{
		Title:           faker.Sentence(30),
		Body:            faker.Paragraph(4, 20, 20, "."),
		UserID:          0, // Set the appropriate user ID here
		LikeCount:       0,
		CommentCount:    0,
		ViewCount:       0,
		ShareCount:      0,
		IsPromoted:      false,
		PromotionExpiry: time.Now(),
		IsFeatured:      false,
		FeaturedExpiry:  time.Now(),
		Visibility:      Models.VisibilityPublic,
	}

	return post
}
