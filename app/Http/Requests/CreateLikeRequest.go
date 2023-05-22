package requests

type CreateLikeRequest struct {
	LikeableID   uint   `json:"likeable_id" validate:"required"`
	LikeableType string `json:"likeable_type" validate:"required"`
}
