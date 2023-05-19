package controllers

import (
	"errors"
	requests "gonga/app/Http/Requests"
	"gonga/app/Models"
	services "gonga/app/Services"
	cloudinary "gonga/packages/Cloudinary"
	"gonga/utils"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func (c PostController) Index(w http.ResponseWriter, r *http.Request) {
	var posts []Models.Post
	result, err := utils.Paginate(r, c.DB, &posts, "User", "Medias", "Mentions.User", "Hashtags")

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Send the list of users in JSON format
	utils.JSONResponse(w, http.StatusOK, result)
}

func (c PostController) Show(w http.ResponseWriter, r *http.Request) {
	// Handle GET /postcontroller/{id} request
	postId, err := utils.GetParam(r, "id")
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch user from the database
	var post Models.Post
	if err := c.DB.Where("id = ?", postId).
		Preload("Medias").
		Preload("Mentions.User").
		Preload("User").
		Preload("Hashtags", func(db *gorm.DB) *gorm.DB {
			// Exclude the "User" field from being loaded for hashtags
			return db.Omit("User")
		}).
		First(&post).Error; err != nil {
		// User not found, return error response
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Return successful response with the user data
	utils.JSONResponse(w, http.StatusOK, post)
}

func (c PostController) Create(w http.ResponseWriter, r *http.Request) {
	// Parse update request from request body
	var createReq requests.CreatePostRequest

	if err := utils.DecodeJSONBody(w, r, &createReq); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			utils.JSONResponse(w, mr.Status(), map[string]string{"error": mr.Error()})
		} else {
			log.Print(err.Error())
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	// Validate post request
	if err := utils.ValidateRequest(w, &createReq); err != nil {
		return
	}
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	newPost := Models.Post{
		Title:           createReq.Title,
		Body:            createReq.Body,
		IsPromoted:      createReq.IsPromoted,
		IsFeatured:      createReq.IsFeatured,
		Visibility:      createReq.Visibility,
		PromotionExpiry: createReq.PromotionExpiry,
		FeaturedExpiry:  createReq.FeaturedExpiry,
		UserID:          uint(userID.(float64)),
	}
	// Insert the post in the database
	result := c.DB.Create(&newPost)

	if result.Error != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		return
	}

	// Associate the media files with the post
	for _, newMedia := range createReq.Medias {
		mediaID := newMedia.ID // Get the ID of the uploaded media

		// Fetch the media record from the database
		var media Models.Media
		if err := c.DB.First(&media, mediaID).Error; err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to associate media with the post"})
			return
		}
		// Update the owner ID of the media to the ID of the newly created post
		media.OwnerID = newPost.ID
		c.DB.Save(&media)
	}

	// Iterate over the mention user IDs
	for _, mentionedUser := range createReq.Mentions {
		log.Println(mentionedUser.UserID)
		mention := &Models.Mention{
			UserID:    mentionedUser.UserID,
			OwnerID:   newPost.ID,
			OwnerType: "posts",
		}
		// Save the mention to the database
		if err := c.DB.Create(&mention).Error; err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
			return
		}
	}

	// Create a slice to store the tags
	var tags []Models.Tag

	// Iterate over the tag titles
	for _, hashtag := range createReq.Hashtags {
		tag := Models.Tag{}

		// Check if the tag already exists in the database
		if err := c.DB.Where("title = ?", hashtag.Title).First(&tag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create a new tag since it doesn't exist
				tag = Models.Tag{
					Title:  hashtag.Title,
					UserID: uint(userID.(float64)),
					// Set other tag fields as needed
				}
				if err := c.DB.Create(&tag).Error; err != nil {
					utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
					return
				}
			} else {
				utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
		}

		// Append the tag to the slice
		tags = append(tags, tag)
	}

	// Associate the tags with the post
	if err := c.DB.Model(&newPost).Association("Hashtags").Replace(tags); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Send email or notification to subscribed users
	// if err := sendPasswordResetEmail(passwordReset.Email, token); err != nil {
	// 	utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	// 	return
	// }
	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Post created successfully!"})
}

func (c PostController) Update(w http.ResponseWriter, r *http.Request) {
	// Handle PUT /postcontroller/{id} request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	defer file.Close()
	// Generate a unique public ID for the file
	publicID := uuid.New().String()

	cloudinaryClient := cloudinary.NewCloudinaryClient()
	result, err := cloudinaryClient.UploadImage(file, publicID)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	utils.JSONResponse(w, http.StatusOK, result)
}

// Update Post Title
func (c PostController) UpdateTitle(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	postID, err := utils.GetParam(r, "id")

	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse update data from request body
	var updateReq requests.UpdatePostTitleRequest
	if err := utils.DecodeJSONBody(w, r, &updateReq); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			utils.JSONResponse(w, mr.Status(), map[string]string{"error": mr.Error()})
		} else {
			log.Print(err.Error())
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	// Validate post request
	if err := utils.ValidateRequest(w, &updateReq); err != nil {
		return
	}

	// Perform update in the database for the specified post ID
	var post Models.Post
	result := c.DB.First(&post, postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		} else {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}

		return
	}

	if post.UserID != uint(userID.(float64)) {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "You are not authorized to update post."})
		return
	}
	// Update the post title
	post.Title = updateReq.Title

	result = c.DB.Save(&post)
	if result.Error != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		return
	}
	// Return success response
	utils.JSONResponse(w, http.StatusOK, map[string]string{"error": "Post title updated successfully!"})
}

// Update Post Body
func (c PostController) UpdateBody(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	postID, err := utils.GetParam(r, "id")

	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse update data from request body
	var updateReq requests.UpdatePostBodyRequest
	if err := utils.DecodeJSONBody(w, r, &updateReq); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			utils.JSONResponse(w, mr.Status(), map[string]string{"error": mr.Error()})
		} else {
			log.Print(err.Error())
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	// Validate post request
	if err := utils.ValidateRequest(w, &updateReq); err != nil {
		return
	}

	// Perform update in the database for the specified post ID
	var post Models.Post
	result := c.DB.First(&post, postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		} else {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}

		return
	}

	if post.UserID != uint(userID.(float64)) {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "You are not authorized to update post."})
		return
	}
	// Update the post title
	post.Body = updateReq.Body

	result = c.DB.Save(&post)
	if result.Error != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		return
	}

	// Perform the edit mentions operation
	err = services.EditMentions(c.DB, post.ID, "posts", updateReq.Mentions)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Return success response
	utils.JSONResponse(w, http.StatusOK, map[string]string{"error": "Post body updated successfully!"})
}

func (c PostController) UpdateMedia(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	postID, err := utils.GetParam(r, "id")

	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse update data from request body
	var updateReq requests.UpdatePostMediaRequest
	if err := utils.DecodeJSONBody(w, r, &updateReq); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			utils.JSONResponse(w, mr.Status(), map[string]string{"error": mr.Error()})
		} else {
			log.Print(err.Error())
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	// Validate post request
	if err := utils.ValidateRequest(w, &updateReq); err != nil {
		return
	}

	// Perform update in the database for the specified post ID
	var post Models.Post
	result := c.DB.First(&post, postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		} else {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}

		return
	}

	if post.UserID != uint(userID.(float64)) {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "You are not authorized to update post."})
		return
	}

	// Perform the edit mentions operation
	err = services.EditMedia(c.DB, post.ID, "posts", updateReq.Medias)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Return success response
	utils.JSONResponse(w, http.StatusOK, map[string]string{"error": "Post medias updated successfully!"})
}

func (c PostController) UpdateHashtag(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	postID, err := utils.GetParam(r, "id")

	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse update data from request body
	var updateReq requests.UpdatePostHashtagRequest
	if err := utils.DecodeJSONBody(w, r, &updateReq); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			utils.JSONResponse(w, mr.Status(), map[string]string{"error": mr.Error()})
		} else {
			log.Print(err.Error())
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	// Validate post request
	if err := utils.ValidateRequest(w, &updateReq); err != nil {
		return
	}

	// Perform update in the service for the specified post ID
	err = services.EditTags(c.DB, postID, updateReq.Hashtags, uint(userID.(float64)))
	if err != nil {
		if err.Error() == "post not found" {
			utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		} else {
			log.Println(err.Error())
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update post tags"})
		}
		return
	}

	// Return success response
	utils.JSONResponse(w, http.StatusOK, map[string]string{"error": "Post medias updated successfully!"})
}

func (c PostController) Delete(w http.ResponseWriter, r *http.Request) {
	// Handle DELETE /postcontroller/{id} request
}
