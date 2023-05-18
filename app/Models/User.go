package Models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username           string     `gorm:"uniqueIndex:idx_username_length;not null;unique"`
	Email              string     `gorm:"unique; not null"`
	Password           string     `gorm:"not null"`
	FirstName          string     `gorm:"not null"`
	LastName           string     `gorm:"not null"`
	AvatarURL          string     `json:"avatar_url"`
	Bio                string     `json:"bio"`
	Gender             string     `json:"gender"`
	MobileNo           string     `json:"mobile_no" gorm:"type:varchar(255);null"`
	MobileNoCode       string     `json:"mobile_no_code"`
	Birthday           *time.Time `json:"birthday"`
	Country            string     `json:"country"`
	City               string     `json:"city"`
	Posts              []Post     `json:"posts" gorm:"foreignKey:UserID"`
	Comments           []Comment  `json:"comments" gorm:"foreignKey:UserID"`
	FollowersList      []Follow   `json:"followers_list" gorm:"foreignKey:FollowerID"`
	FollowingList      []Follow   `json:"following_list" gorm:"foreignKey:FollowingID"`
	BackgroundImageURL string     `json:"background_image_url"`
	WebsiteURL         string     `json:"website_url"`
	Occupation         string     `json:"occupation"`
	Education          string     `json:"education"`
	// Interests          []string  `json:"interests"`
}

func (User) TableName() string {
	return "users"
}
