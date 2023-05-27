// factory.go

package factory

import (
	"gonga/app/Models"

	"github.com/bxcodec/faker/v3"
)

// GenerateUser generates a fake user instance
func GenerateUser() Models.User {
	user := Models.User{
		FirstName:     faker.FirstName(),
		LastName:      faker.LastName(),
		Email:         faker.Email(),
		Gender:        faker.Gender(),
		Username:      faker.Username(),
		Bio:           faker.Paragraph(),
		AvatarURL:     faker.URL(),
		WebsiteURL:    faker.URL(),
		MobileNo:      faker.Phonenumber(),
		EmailVerified: true,
		Password:      faker.Password(), // You can set a default password here
	}

	return user
}
