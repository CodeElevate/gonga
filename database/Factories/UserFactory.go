// factory.go

package factory

import (
	"gonga/app/Models"

	"github.com/go-faker/faker/v4"
)

// GenerateUser generates a fake user instance
func UserFactory() Models.User {
	user := Models.User{
		Username:     faker.Username(),
		Email:        faker.Email(),
		Password:     faker.Password(),
		FirstName:    faker.FirstName(),
		LastName:     faker.LastName(),
		AvatarURL:    faker.URL(),
		Bio:          faker.Sentence(),
		Gender:       faker.Gender(),
		MobileNo:     faker.Phonenumber(),
		MobileNoCode: "+1",
		// Birthday:    faker.DateTime.Date(),
		// Country:     faker,
		// City:        faker.City(),
		WebsiteURL: faker.URL(),
		// Occupation:  faker.JobTitle(),
		Education:     faker.Word(),
		EmailVerified: true,
	}

	return user
}
