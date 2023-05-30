package seeder

import (
	"fmt"
	"gonga/database/Factories"
	"gorm.io/gorm"
	"log"
)

type UserSeeder struct {
	DB *gorm.DB
}

func (s *UserSeeder) Run() {
	// Access the database connection
	db := s.DB

	// Create and save dummy user records using the factory
	for i := 0; i < 100; i++ {
		user := factory.UserFactory()
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Error seeding user: %v", err)
		}
	}

	fmt.Println("User seeding completed.")
}
