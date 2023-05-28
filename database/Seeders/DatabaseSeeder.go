package seeder

import (
	"gorm.io/gorm"
)

type DatabaseSeeder struct {
	DB *gorm.DB
}

func (s *DatabaseSeeder) Run() {
	// Instantiate individual seeders
	userSeeder := &UserSeeder{DB: s.DB}
	tagSeeder := &TagSeeder{DB: s.DB}
	postSeeder := &PostSeeder{DB: s.DB}
	commentSeeder := &CommentSeeder{DB: s.DB}

	// Run the seeders
	userSeeder.Run()
	tagSeeder.Run()
	postSeeder.Run()
	commentSeeder.Run()
}
