package commands

import (
	"fmt"
	"gonga/bootstrap"
	seeder "gonga/database/Seeders"

	"github.com/spf13/cobra"
)

func SeedCmd(app *bootstrap.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "db:seed",
		Short: "Run database migrations.",
		Long:  "Run any pending database migrations.",
		Run: func(_ *cobra.Command, _ []string) {
			// Open a connection to the database using the configured credentials
			db := app.DB

			// Create an instance of each seeder
			userSeeder := &seeder.UserSeeder{DB: db}

			// Run the seeders
			userSeeder.Run()

			// Add more seeders as needed

			// Print completion message
			fmt.Println("Database seeding completed.")
		},
	}
}
