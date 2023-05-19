package requests

import "gonga/app/Models"

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
