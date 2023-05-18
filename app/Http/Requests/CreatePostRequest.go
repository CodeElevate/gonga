package requests

import (
	"gonga/app/Models"
	"time"
)

type CreatePostRequest struct {
	Title string `json:"title" validate:"required,min=20"`
	Body  string `json:"body" validate:"required,min=40"`
	Hashtags        []Models.Tag   `json:"hashtags" validate:"min=1,max=5"`
	Mentions        []Models.User  `json:"mentions" validate:"omitempty,max=5"`
	Medias          []Models.Media `json:"medias" validate:"omitempty,max=5"`
	IsPromoted      bool      `json:"is_promoted"`
	IsPublished     bool      `json:"is_published"`
	PromotionExpiry time.Time `json:"promotion_expiry"`
	IsFeatured      bool      `json:"is_featured"`
	FeaturedExpiry  time.Time `json:"featured_expiry"`
}

// UserID          uint      `json:"user_id"`
//     User            User      `json:"user"`
