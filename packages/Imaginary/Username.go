package imaginary

import (
	"strconv"
	"strings"
	"sync"

	faker "github.com/brianvoe/gofakeit/v6"
)

// UserNameGenerator generates unique usernames
type UserNameGenerator struct {
	usedUserNames map[string]bool
	mutex         sync.Mutex
}

// NewUserNameGenerator creates a new UserNameGenerator instance
func NewUserNameGenerator() *UserNameGenerator {
	return &UserNameGenerator{
		usedUserNames: make(map[string]bool),
	}
}

// GenerateUniqueUserName generates a case-insensitive unique username
func (g *UserNameGenerator) UserName() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	baseUserName := faker.Username()
	userName := baseUserName
	counter := 1

	for g.isUserNameTaken(userName) {
		userName = baseUserName + "_" + strconv.Itoa(counter)
		counter++
	}

	g.usedUserNames[strings.ToLower(userName)] = true
	return userName
}

// isUserNameTaken checks if a username is already taken (case-insensitive)
func (g *UserNameGenerator) isUserNameTaken(userName string) bool {
	for usedName := range g.usedUserNames {
		if strings.EqualFold(usedName, userName) {
			return true
		}
	}
	return false
}
