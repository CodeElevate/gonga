package seeder

import (
	"fmt"
	"gonga/app/Models"
	factory "gonga/database/Factories"
	imaginary "gonga/packages/Imaginary"
	"log"

	"gorm.io/gorm"
)

type TagSeeder struct {
	DB *gorm.DB
}

func (s *TagSeeder) Run() {
	// Access the database connection
	db := s.DB

	// Retrieve some users from the database
	var users []Models.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("Error retrieving users: %v", err)
	}

	// Create and save dummy tag records using the factory
	tagGenerator := imaginary.NewTagGenerator()
	for i := 0; i < 100; i++ {
		tag := factory.TagFactory()
		tag.User = users[i%len(users)]
		tag.Title = tagGenerator.UniqueTag()
		if err := db.Create(&tag).Error; err != nil {
			log.Fatalf("Error seeding tag: %v", err)
		}
	}

	fmt.Println("Tag seeding completed.")
}
