package requests

import (
	"gonga/app/Models"
	"time"
)

type UpdatePostRequest struct {
}

type UpdatePostTitleRequest struct {
	Title string `json:"title" validate:"required,min=20"`
}

type UpdatePostBodyRequest struct {
	Body     string           `json:"body" validate:"required,min=40"`
	Mentions []Models.Mention `json:"mentions" validate:"omitempty,max=15"`
}

type UpdatePostMediaRequest struct {
	Medias []Models.Media `json:"medias" validate:"omitempty,max=15"`
}

type UpdatePostHashtagRequest struct {
	Hashtags []Models.Tag `json:"hashtags" validate:"min=1,max=5"`
}

type UpdatePostSettingsRequest struct {
	IsPromoted      *bool              `json:"is_promoted" validate:"required"`
	PromotionExpiry time.Time         `json:"promotion_expiry" validate:"required"`
	IsFeatured      *bool              `json:"is_featured" validate:"required"`
	FeaturedExpiry  *time.Time         `json:"featured_expiry" validate:"required"`
	Visibility      Models.Visibility `json:"visibility" validate:"required"`
}
