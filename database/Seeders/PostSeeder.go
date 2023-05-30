package seeder

import (
	"fmt"
	"gonga/app/Models"
	factory "gonga/database/Factories"
	"log"

	"gorm.io/gorm"
)

type PostSeeder struct {
	DB *gorm.DB
}

func (s *PostSeeder) Run() {
	// Access the database connection
	db := s.DB

	// Retrieve some users from the database
	var users []Models.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("Error retrieving users: %v", err)
	}

	// Retrieve some tags from the database
	var tags []Models.Tag
	if err := db.Find(&tags).Error; err != nil {
		log.Fatalf("Error retrieving tags: %v", err)
	}

	// Create and save dummy post records using the factory
	for i := 0; i < 200; i++ {
		post := factory.PostFactory()
		post.User = &users[i%len(users)]
		post.Hashtags = append(post.Hashtags, &tags[i%len(tags)])

		if err := db.Create(&post).Error; err != nil {
			log.Fatalf("Error seeding post: %v", err)
		}
	}

	fmt.Println("Post seeding completed.")
}
