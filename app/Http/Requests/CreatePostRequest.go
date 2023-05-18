package requests

import (
	"gonga/app/Models"
	"time"
)

type CreatePostRequest struct {
	Title           string            `json:"title" validate:"required,min=20"`
	Body            string            `json:"body" validate:"required,min=40"`
	Hashtags        []Models.Tag      `json:"hashtags" validate:"min=1,max=5"`
	Mentions        []Models.Mention  `json:"mentions" validate:"omitempty,max=15"`
	Medias          []Models.Media    `json:"medias" validate:"omitempty,max=15"`
	IsPromoted      bool              `json:"is_promoted"`
	PromotionExpiry time.Time         `json:"promotion_expiry"`
	IsFeatured      bool              `json:"is_featured"`
	FeaturedExpiry  time.Time         `json:"featured_expiry"`
	Visibility      Models.Visibility `json:"visibility"`
}
