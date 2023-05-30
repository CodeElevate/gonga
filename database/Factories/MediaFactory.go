package factory

import (
	"gonga/app/Models"
	imaginary "gonga/packages/Imaginary"

	"github.com/jaswdr/faker"
)

// MediaFactory generates a fake media instance
func MediaFactory() Models.Media {
	faker := faker.New()
	image := imaginary.GenerateImage(800, 600)

	media := Models.Media{
		URL:       image.URL,
		Type:      image.Type,
		OwnerID:   0,                                                                 // Set the appropriate owner ID here
		OwnerType: faker.RandomStringElement([]string{"posts", "comments", "users"}), // Set the appropriate owner type here
	}

	return media
}
