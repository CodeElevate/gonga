package factory

import (
	"gonga/app/Models"

	faker "github.com/brianvoe/gofakeit/v6"
)

// MediaFactory generates a fake media instance
func MediaFactory() Models.Media {

	media := Models.Media{
		URL:       faker.ImageURL(300, 400),
		Type:      faker.RandomString([]string{"jpeg", "png", "gif"}),
		OwnerID:   0,                                                          // Set the appropriate owner ID here
		OwnerType: faker.RandomString([]string{"posts", "comments", "users"}), // Set the appropriate owner type here
	}

	return media
}
