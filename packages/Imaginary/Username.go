package imaginary

import (
	"github.com/jaswdr/faker"
	"sync"
)
// UserNameGenerator generates unique user names
type UserNameGenerator struct {
	faker        *faker.Faker
	usedUserNames map[string]bool
	mutex        sync.Mutex
}

// NewUserNameGenerator creates a new UserNameGenerator instance with a pre-initialized faker.Faker instance
func NewUserNameGenerator(faker *faker.Faker) *UserNameGenerator {
	return &UserNameGenerator{
		faker:        faker,
		usedUserNames: make(map[string]bool),
	}
}

// GenerateUniqueUserName generates a unique user name
func (g *UserNameGenerator) UserName() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	userName := g.faker.Person().Name()
	for g.usedUserNames[userName] {
		userName = g.faker.Person().Name()
	}

	g.usedUserNames[userName] = true
	return userName
}
