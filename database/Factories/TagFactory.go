package factory

import (
	"gonga/app/Models"
	imaginary "gonga/packages/Imaginary"

	"github.com/jaswdr/faker"
)

// TagFactory generates a fake tag instance
func TagFactory() Models.Tag {
	faker := faker.New()
	tag := Models.Tag{
		// Title:        faker.Username(),
		CoverImage:   imaginary.GenerateImage(400, 400).URL,
		BackendImage: imaginary.GenerateImage(800, 600).URL,
		Description:  faker.Lorem().Paragraph(10),
		Color:        faker.Color().Hex(),
		Slug:         faker.Internet().Slug(),
		UserID:       0, // Set the appropriate user ID here
	}

	return tag
}
