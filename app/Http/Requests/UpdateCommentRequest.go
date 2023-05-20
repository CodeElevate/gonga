package requests

import "gonga/app/Models"

type UpdateCommentRequest struct {
	Body     string           `json:"body" validate:"required,min=40"`
	Mentions []Models.Mention `json:"mentions" validate:"omitempty,max=15"`
}