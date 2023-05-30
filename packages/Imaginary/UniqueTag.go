package imaginary

import (
	"strconv"
	"sync"

	faker "github.com/brianvoe/gofakeit/v6"
)

// TagGenerator generates unique tags
type TagGenerator struct {
	usedTags map[string]bool
	mutex    sync.Mutex
}

// NewTagGenerator creates a new TagGenerator instance
func NewTagGenerator() *TagGenerator {
	return &TagGenerator{
		usedTags: make(map[string]bool),
	}
}

// GenerateUniqueTag generates a unique tag
func (g *TagGenerator) UniqueTag() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	baseTag := faker.Word()
	tag := baseTag
	counter := 1

	for g.usedTags[tag] {
		tag = baseTag + "_" + strconv.Itoa(counter)
		counter++
	}

	g.usedTags[tag] = true
	return tag
}
