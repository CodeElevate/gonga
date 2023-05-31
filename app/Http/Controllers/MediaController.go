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

//	@Summary		Upload media files
//	@Description	Upload media files
//	@Tags			Media
//	@Accept			mpfd
//	@Produce		json
//	@Param			files		formData	file							true	"Media files to upload"
//	@Param			owner_type	formData	string							false	"Type of the owner of the media file"
//	@Param			owner_id	formData	string							false	"ID of the owner of the media file"
//	@Success		200			{array}		responses.UploadMediaResponse	"success"
//	@Failure		400			{object}	utils.APIResponse				"error"
//	@Router			/upload [POST]
//	@Security		BearerAuth
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
		ownerType = "posts" // Fallback value for owner type (assuming it's a post)
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
