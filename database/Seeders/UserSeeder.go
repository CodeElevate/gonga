package seeder

import (
	"fmt"
	factory "gonga/database/Factories"
	imaginary "gonga/packages/Imaginary"
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
	UserNameGenerator := imaginary.NewUserNameGenerator()
	for i := 0; i < 500; i++ {
		user := factory.UserFactory()
		user.Username = UserNameGenerator.UserName()
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Error seeding user: %v", err)
		}
		if i%1000 == 0 && i > 0 {
			fmt.Println(i, "new user added")
		}
	}

	fmt.Println("User seeding completed.")
}
