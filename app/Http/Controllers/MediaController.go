package controllers

import (
	responses "gonga/app/Http/Responses"
	"gonga/app/Models"
	cloudinary "gonga/packages/Cloudinary"
	"gonga/utils"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaController struct {
	DB *gorm.DB
}

// Upload handles the POST /upload request and uploads multiple files to the server.
// It parses the multipart form, retrieves the files from the request, and uploads them to a cloud storage service (Cloudinary).
// The uploaded media records are also stored in the database.
//
// Parameters:
// - w: http.ResponseWriter - the response writer used to send HTTP responses.
// - r: *http.Request - the HTTP request containing the uploaded files.
//
// Response:
// The function sends an HTTP response with the uploaded media information in JSON format.
// The response includes the URL, type, filename, size, and ID for each uploaded file.
//
// Example:
//
//	POST /upload HTTP/1.1
//	Host: example.com
//	Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
//
//	------WebKitFormBoundary7MA4YWxkTrZu0gW
//	Content-Disposition: form-data; name="files"; filename="image1.jpg"
//	Content-Type: image/jpeg
//
//	<binary data>
//
//	------WebKitFormBoundary7MA4YWxkTrZu0gW
//	Content-Disposition: form-data; name="files"; filename="image2.jpg"
//	Content-Type: image/jpeg
//
//	<binary data>
//
//	------WebKitFormBoundary7MA4YWxkTrZu0gW--
func (c MediaController) Upload(w http.ResponseWriter, r *http.Request) {
	// Handle POST /upload request
	err := r.ParseMultipartForm(32 << 20) // Max upload file size: 32MB
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to parse multipart form"})
		return
	}

	files := r.MultipartForm.File["files"]      // Assuming "files" is the name of the file input field for multiple uploads
	var results []responses.UploadMediaResponse // Store the upload results for each file

	ownerType := r.FormValue("owner_type") // Get the owner type from the request form data
	if ownerType == "" {
		ownerType = "post" // Fallback value for owner type (assuming it's a post)
	}
	ownerID := r.FormValue("owner_id") // Get the owner ID from the request form data

	cloudinaryClient := cloudinary.NewCloudinaryClient()

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		defer file.Close()

		// Generate a unique public ID for the file
		publicID := uuid.New().String()

		result, err := cloudinaryClient.UploadImage(file, publicID)
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		// Store the media record in the database
		media := Models.Media{
			URL:       result.URL,
			Type:      result.Type,           // Set the appropriate media type
			OwnerType: ownerType,             // Set the owner type dynamically or fallback to "post"
			OwnerID:   parseOwnerID(ownerID), // Parse the owner ID based on its type (post, comment, etc.)

		}
		c.DB.Create(&media)
		// Send response
		response := responses.UploadMediaResponse{
			URL:      result.URL,
			Type:     result.Type,
			Filename: result.Filename,
			Size:     result.Size,
			ID:       media.ID,
		}
		results = append(results, response)
	}

	utils.JSONResponse(w, http.StatusOK, results)
}

// parseOwnerID is a helper function to parse the owner ID based on its type.
// Modify this function based on your actual implementation.
func parseOwnerID(ownerID string) uint {
	// Parse the owner ID based on its type (post, comment, etc.)
	// Modify this function based on your actual implementation.
	// Example: Convert ownerID to uint and return it
	parsedOwnerID, _ := strconv.ParseUint(ownerID, 10, 64)
	return uint(parsedOwnerID)
}
