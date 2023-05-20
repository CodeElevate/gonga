package requests

import "gonga/app/Models"

type CreateCommentRequest struct {
	Body     string         `json:"body" validate:"required"`
	ParentID *uint           `json:"parent_id,omitempty"`
	Mentions        []Models.Mention  `json:"mentions" validate:"omitempty,max=15"`
}
