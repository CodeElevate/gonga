package Models

import (
    "gorm.io/gorm"
)

type Notification struct {
    gorm.Model
    Name    string   `gorm:"not null"`
    Email   string   `gorm:"unique;not null"`
}

func (Notification) TableName() string {
    return "notifications"
}