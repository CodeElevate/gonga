package seeder

import (
	"fmt"
	factory "gonga/database/Factories"
	"log"

	"gorm.io/gorm"
)

type UserSeeder struct {
	DB *gorm.DB
}

func (s *UserSeeder) Run() {
	// Access the database connection
	db := s.DB

	// Create and save dummy user records using the factory
	for i := 0; i < 10; i++ {
		user := factory.GenerateUser()
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Error seeding user: %v", err)
		}
	}

	fmt.Println("User seeding completed.")
}
