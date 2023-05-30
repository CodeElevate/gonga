package factory

import (
	"gonga/app/Models"
	"time"

	"github.com/jaswdr/faker"
)

// PostFactory generates a fake post instance
func PostFactory() Models.Post {
	faker := faker.New()

	post := Models.Post{
		Title:           faker.Lorem().Sentence(30),
		Body:            faker.Lorem().Paragraph(4),
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
