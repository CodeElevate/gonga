// factory.go

package factory

import (
	"gonga/app/Models"
	"gonga/utils"
	"log"

	faker "github.com/brianvoe/gofakeit/v6"
)

// GenerateUser generates a fake user instance
func UserFactory() Models.User {
	birthday := faker.Date()
	// Create user in database
	hashedPassword, err := utils.HashPassword(faker.Password(true, true, true, true, false, 12))

	if err != nil {
		log.Println(err.Error())
	}
	// imaginary.NewUserNameGenerator().UserName()
	user := Models.User{
		Email:         faker.Email(),
		Password:      hashedPassword,
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
