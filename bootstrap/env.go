package bootstrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
)

func LoadEnv() {
	_, err := os.Stat(".env")
	if err == nil {
		// If the .env file exists, load it
		err = godotenv.Load(".env")
		if err != nil {
			pterm.Fatal.Printf("Error loading .env file: %s", err)
		}
	} else if os.IsNotExist(err) {
		// If the .env file doesn't exist, rely on built-in environment variables (e.g., Railway)
		log.Println(".env file not found. Relying on built-in environment variables.")
	} else {
		// Handle any other error that occurred while checking the file status
		pterm.Fatal.Printf("Error checking .env file: %s", err)
	}
}
