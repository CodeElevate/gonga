package factory

import (
	"gonga/app/Models"
	"github.com/go-faker/faker/v4"
)

// TagFactory generates a fake tag instance
func TagFactory() Models.Tag {
	tag := Models.Tag{
		Title:        faker.Word(),
		CoverImage:   faker.URL(),
		BackendImage: faker.URL(),
		Description:  faker.Sentence(),
		// Color:        faker.Color().Hex(),
		Slug:         faker.Username(),
		UserID:       0, // Set the appropriate user ID here
	}

	return tag
}
