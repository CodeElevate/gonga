package Models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID          uint      `json:"user_id"`
    User            User      `json:"user"`
	Title           string    `json:"title"`
	Body            string    `json:"body"`
	Likes           []Like    `json:"likes" gorm:"foreignKey:PostID"`
	LikeCount       uint      `json:"likeCount"`
	Comments        []Comment `json:"comments" gorm:"foreignKey:PostID"`
	CommentCount    uint      `json:"comment_count"`
	ViewCount       uint      `json:"view_count"`
	ShareCount      uint      `json:"share_count"`
	Medias          []*Media  `json:"medias" gorm:"polymorphic:Owner;"`
	Hashtags        []*Tag    `json:"hashtags" gorm:"many2many:post_hashtags;"`
	Mentions        []*User   `json:"mentions" gorm:"many2many:post_mentions;"`
	IsPromoted      bool      `json:"is_promoted"`
	IsPublished      bool      `json:"is_published"`
	PromotionExpiry time.Time `json:"promotion_expiry"`
	IsFeatured      bool      `json:"is_featured"`
	FeaturedExpiry  time.Time `json:"featured_expiry"`
}

func (Post) TableName() string {
	return "posts"
}
