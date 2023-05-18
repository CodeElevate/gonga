package Models

import (
    "gorm.io/gorm"
)

type Media struct {
    gorm.Model
    URL       string `json:"url"`
    Type      string `json:"type"`
    OwnerID   uint   `json:"owner_id"`
    OwnerType string `json:"owner_type"` // post, comment, profile, etc.
}

func (Media) TableName() string {
    return "medias"
}