package Models

import (
    "gorm.io/gorm"
)

type Media struct {
    gorm.Model
    URL       string `json:"url"`
    Type      string `json:"type"`
    OwnerID   uint   `json:"owner_id"`
    OwnerType string `json:"owner_type"` // posts, comments, users, etc.
}

func (Media) TableName() string {
    return "medias"
}