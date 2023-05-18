package Models

import (
    "gorm.io/gorm"
)

type Media struct {
    gorm.Model
    URL         string `json:"url"`
    Type        string `json:"type"`
    OwnerType   string `json:"owner_type"` // post, comment, profile, etc.
    OwnerID     uint   `json:"owner_id"`
}

func (Media) TableName() string {
    return "medias"
}