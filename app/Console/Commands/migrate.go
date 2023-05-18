package commands

import (
	"fmt"
	"gonga/app/Models"
	"gonga/bootstrap"
	"log"

	"github.com/spf13/cobra"
)

func MigrateCmd(app *bootstrap.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations.",
		Long:  "Run any pending database migrations.",
		Run: func(_ *cobra.Command, _ []string) {
			// Open a connection to the database using the configured credentials
			db := app.DB
			// Register the migrations
			db.Debug()
			err := db.AutoMigrate(
				&Models.Comment{},
				&Models.Follow{},
				&Models.Like{},
				&Models.Media{},
				&Models.Tag{},
				&Models.User{},
				&Models.Post{},
				&Models.PasswordReset{},
				&Models.PersonalAccessToken{},
				&Models.Mention{},
			)
			if err != nil {
				log.Fatalf("Error running migrations: %v", err)
			}

			// Run any pending database migrations using the configured migration tool (e.g. gorm)
			// ...

			fmt.Println("Database migrations completed.")
		},
	}
}
