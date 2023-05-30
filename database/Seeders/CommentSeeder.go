package seeder

import (
	"fmt"
	"gonga/app/Models"
	factory "gonga/database/Factories"
	"log"

	"gorm.io/gorm"
)

type CommentSeeder struct {
	DB *gorm.DB
}

func (s *CommentSeeder) Run() {
	// Access the database connection
	db := s.DB

	// Retrieve some users from the database
	var users []Models.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("Error retrieving users: %v", err)
	}

	// Retrieve some posts from the database
	var posts []Models.Post
	if err := db.Find(&posts).Error; err != nil {
		log.Fatalf("Error retrieving posts: %v", err)
	}

	// Create and save dummy comment records using the factory
	for i := 0; i < 400000; i++ {
		comment := factory.CommentFactory()
		comment.User = &users[i%len(users)]
		comment.Post = &posts[i%len(posts)]

		if err := db.Create(&comment).Error; err != nil {
			log.Fatalf("Error seeding comment: %v", err)
		}
		if i%5000 == 0 && i > 0 {
			fmt.Println(i, "new comments added")
		}
	}

	fmt.Println("Comment seeding completed.")
}
