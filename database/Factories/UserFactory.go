// factory.go

package factory

import (
	"gonga/app/Models"

	faker "github.com/brianvoe/gofakeit/v6"
)

// GenerateUser generates a fake user instance
func UserFactory() Models.User {
	birthday := faker.Date()
	user := Models.User{
		Username:      faker.Username(),
		Email:         faker.Email(),
		Password:      faker.Password(true, true, true, true, false, 12),
		FirstName:     faker.FirstName(),
		LastName:      faker.LastName(),
		AvatarURL:     faker.Person().Image,
		Bio:           faker.Paragraph(2, 20, 20, "."),
		Gender:        faker.Gender(),
		MobileNo:      faker.Phone(),
		MobileNoCode:  "+91",
		Birthday:      &birthday,
		Country:       faker.Country(),
		City:          faker.City(),
		WebsiteURL:    faker.URL(),
		Occupation:    faker.JobTitle(),
		Education:     faker.BS(),
		EmailVerified: true,
	}

	return user
}
