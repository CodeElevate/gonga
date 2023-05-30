// factory.go

package factory

import (
	"gonga/app/Models"
	"time"

	"github.com/jaswdr/faker"
)

// GenerateUser generates a fake user instance
func UserFactory() Models.User {
	faker := faker.New()

	user := Models.User{
		Username:      faker.Username(),
		Email:         faker.Internet().SafeEmail(),
		Password:      faker.Internet().Password(),
		FirstName:     faker.Person().FirstName(),
		LastName:      faker.Person().LastName(),
		AvatarURL:     faker.URL(),
		Bio:           faker.Lorem().Sentence(400),
		Gender:        faker.Person().Gender(),
		MobileNo:      faker.Phone().Number(),
		MobileNoCode:  faker.Address().CountryCode(),
		Birthday:      &time.Time{},
		Country:       faker.Address().Country(),
		City:          faker.Address().City(),
		WebsiteURL:    faker.Internet().URL(),
		Occupation:    faker.Company().JobTitle(),
		Education:     faker.Person().Title(),
		EmailVerified: true,
	}

	return user
}
