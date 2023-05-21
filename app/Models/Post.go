package Models

import (
	"time"

	"gorm.io/gorm"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
	VisibilityFriends Visibility = "friends"
	// Add more visibility options as needed
)

type Post struct {
	gorm.Model
	UserID          uint       `json:"user_id"`
	User            *User      `json:"user,omitempty"`
	Title           string     `json:"title"`
	Body            string     `json:"body"`
	Likes           []Like     `json:"likes" gorm:"polymorphic:Likeable;"`
	LikeCount       uint       `json:"like_count"`
	Comments        []Comment  `json:"comments" gorm:"foreignKey:PostID"`
	CommentCount    uint       `json:"comment_count"`
	ViewCount       uint       `json:"view_count"`
	ShareCount      uint       `json:"share_count"`
	Medias          []*Media   `json:"medias" gorm:"polymorphic:Owner;"`
	Hashtags        []*Tag     `json:"hashtags" gorm:"many2many:post_hashtags;"`
	Mentions        []*Mention `json:"mentions" gorm:"polymorphic:Owner;"`
	IsPromoted      bool       `json:"is_promoted"`
	PromotionExpiry time.Time  `json:"promotion_expiry"`
	IsFeatured      bool       `json:"is_featured"`
	FeaturedExpiry  time.Time  `json:"featured_expiry"`
	Visibility      Visibility `json:"visibility"`
}

func (Post) TableName() string {
	return "posts"
}
