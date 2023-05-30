package factory

import (
	"gonga/app/Models"

	"github.com/jaswdr/faker"
)

// TagFactory generates a fake tag instance
func TagFactory() Models.Tag {
	faker := faker.New()

	tag := Models.Tag{
		Title:        faker.Username(),
		CoverImage:   faker.URL(),
		BackendImage: faker.URL(),
		Description:  faker.Lorem().Paragraph(10),
		Color:        faker.Color().Hex(),
		Slug:         faker.Internet().Slug(),
		UserID:       0, // Set the appropriate user ID here
	}

	return tag
}
