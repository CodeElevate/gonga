package factory

import (
	"gonga/app/Models"

	"github.com/jaswdr/faker"
)

// MediaFactory generates a fake media instance
func MediaFactory() Models.Media {
	faker := faker.New()

	media := Models.Media{
		// URL:       faker.LoremFlickr().,
		Type:      faker.MimeType().MimeType(),
		OwnerID:   0,                                                                 // Set the appropriate owner ID here
		OwnerType: faker.RandomStringElement([]string{"posts", "comments", "users"}), // Set the appropriate owner type here
	}

	return media
}
