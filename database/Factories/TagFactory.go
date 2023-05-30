package factory

import (
	"gonga/app/Models"
	imaginary "gonga/packages/Imaginary"

	faker "github.com/brianvoe/gofakeit/v6"
)

// TagFactory generates a fake tag instance
func TagFactory() Models.Tag {
	tag := Models.Tag{
		Title:        imaginary.NewTagGenerator().UniqueTag(),
		CoverImage:   faker.ImageURL(300, 400),
		BackendImage: faker.ImageURL(800, 600),
		Description:  faker.Paragraph(1, 5, 15, "."),
		Color:        faker.HexColor(),
		Slug:         faker.Sentence(10),
		UserID:       0, // Set the appropriate user ID here
	}

	return tag
}
