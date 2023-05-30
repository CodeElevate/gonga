package imaginary

import (
	"fmt"
	"math/rand"
)

// Image represents a fake image with URL, type, and size information
type Image struct {
	URL  string `json:"url"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

// GenerateURL generates a fake image URL with the specified width and height
func GenerateURL(width, height int) string {
	// Replace "example.com" with your desired base URL or domain name
	return fmt.Sprintf("https://loremflickr.com/%dx%d", width, height)
}

// GenerateType generates a random image type
func GenerateType() string {
	types := []string{"jpeg", "png", "gif"}
	// Select a random type from the available types
	return types[rand.Intn(len(types))]
}

// GenerateSize generates a random image size in kilobytes
func GenerateSize() int {
	// Generate a random size between 100 and 5000 kilobytes
	return rand.Intn(4900) + 100
}

// GenerateImage generates a fake image with URL, type, and size information
func GenerateImage(width, height int) Image {
	url := GenerateURL(width, height)
	imageType := GenerateType()
	size := GenerateSize()

	image := Image{
		URL:  url,
		Type: imageType,
		Size: size,
	}

	return image
}
