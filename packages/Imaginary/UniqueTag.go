package imaginary

import (
	"sync"

	faker "github.com/brianvoe/gofakeit/v6"
)

// TagGenerator generates unique tags
type TagGenerator struct {
	usedTags map[string]bool
	mutex    sync.Mutex
}

// NewTagGenerator creates a new TagGenerator instance with a pre-initialized faker.Faker instance
func NewTagGenerator() *TagGenerator {
	return &TagGenerator{
		usedTags: make(map[string]bool),
	}
}

// GenerateUniqueTag generates a unique tag
func (g *TagGenerator) UniqueTag() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	tag := faker.Word()
	for g.usedTags[tag] {
		tag = faker.Word()
	}

	g.usedTags[tag] = true
	return tag
}
