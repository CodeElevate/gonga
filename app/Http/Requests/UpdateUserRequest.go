package requests

import "time"

// UpdateUserRequest represents the request payload for updating a user's profile information
type UpdateUserRequest struct {
	FirstName          string    `json:"first_name" validate:"required,max=50"`
	LastName           string    `json:"last_name"`
	AvatarURL          string    `json:"avatar_url"`
	Bio                string    `json:"bio" validate:"max=500"`
	Gender             string    `json:"gender" validate:"oneof=male female"`
	MobileNo           string    `json:"mobile_no" validate:"omitempty,numeric,len=10"`
	Country            string    `json:"country" validate:"omitempty,alpha,max=50"`
	City               string    `json:"city" validate:"omitempty,alpha,max=50"`
	MobileNoCode       string    `json:"mobile_no_code" validate:"required"`
	Birthday           *time.Time `json:"birthday" validate:"omitempty,time"`
	BackgroundImageURL string    `json:"background_image_url"`
	WebsiteURL         string    `json:"website_url" validate:"omitempty,url"`
	Occupation         string    `json:"occupation" validate:"omitempty,max=50"`
	Education          string    `json:"education" validate:"omitempty,max=50"`
}
