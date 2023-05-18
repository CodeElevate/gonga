package requests

import "mime/multipart"

type MediaUploadRequest struct {
	File      *multipart.FileHeader `json:"file"`
	OwnerType string                `json:"owner_type"`
	OwnerID   uint                  `json:"owner_id"`
}
