package controllers

import (
	"errors"
	requests "gonga/app/Http/Requests"
	"gonga/app/Models"
	services "gonga/app/Services"
	"gonga/utils"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

// Index retrieves a list of all posts from the server.
//
// This endpoint allows you to fetch all the posts available in the system.
// The returned list includes details such as post ID, title, content, author, and creation date.
//
// The API response is in JSON format.
//
// @Summary Get all posts
// @Description Retrieves a list of all posts from the server.
// @Tags Posts
// @Produce json
// @Success 200 {object} utils.SwaggerPagination
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Router /posts [get]
func (c PostController) Index(w http.ResponseWriter, r *http.Request) {
	var posts []Models.Post
	var response utils.APIResponse

	paginationScope, err := utils.Paginate(r, c.DB, &posts, &response, "User", "Medias", "Mentions.User", "Hashtags")
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "Failed to paginate posts")
		return
	}

	db := paginationScope(c.DB)
	if err := db.Find(&posts).Error; err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	response.Data = posts
	response.Type = "success"
	response.Message = "data retrieved successfully"

	utils.JSONResponse(w, http.StatusOK, response)
}

// Show handles the GET /posts/{id} request to retrieve a specific post.
//
// This endpoint allows users to retrieve the details of a specific post identified by its ID.
//
// @Summary Get a specific post
// @Description Retrieves the details of a specific post
// @Tags Posts
// @Param id path int true "Post ID"
// @Produce json
// @Success 200 {object} utils.SwaggerSuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Router /posts/{id} [get]
func (c PostController) Show(w http.ResponseWriter, r *http.Request) {
	// Handle GET /postcontroller/{id} request
	postId, err := utils.GetParam(r, "id")
	if err != nil {
		log.Println(err.Error())
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	var post Models.Post
	if err := c.DB.Where("id = ?", postId).
		Preload("Medias").
		Preload("Mentions.User").
		Preload("User").
		Preload("Likes").
		Preload("Hashtags", func(db *gorm.DB) *gorm.DB {
			// Exclude the "User" field from being loaded for hashtags
			return db.Omit("User")
		}).
		First(&post).Error; err != nil {
		// User not found, return error response
		utils.HandleError(w, err, http.StatusNotFound)
		return
	}
	// Return successful response with the post data
	response := utils.APIResponse{
		Type: "success",
		Data: post,
	}
	utils.JSONResponse(w, http.StatusOK, response)
}

// Create handles the POST /posts request to create a new post.
//
// This endpoint allows users to create a new post by providing the necessary details in the request body.
// The request body should contain the post title, body, visibility, promotion and featured settings, media files, mentions, and hashtags.
//
// @Summary Create a new post
// @Description Creates a new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param post body requests.CreatePostRequest true "Create post request body"
// @Success 201 {object} utils.SwaggerSuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Router /posts [post]
func (c PostController) Create(w http.ResponseWriter, r *http.Request) {
	// Parse update request from request body
	var createReq requests.CreatePostRequest

	if err := utils.DecodeJSONBody(w, r, &createReq); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			utils.HandleError(w, err, mr.Status())
		} else {
			log.Print(err.Error())
			utils.HandleError(w, errors.New("failed to decode JSON body"), http.StatusInternalServerError)
		}
		return
	}
	// Validate post request
	if err := utils.ValidateRequest(w, &createReq); err != nil {
		return
	}

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
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
		utils.HandleError(w, result.Error, http.StatusInternalServerError, "failed to create post in the database")
		return
	}

	// Associate the media files with the post
	for _, newMedia := range createReq.Medias {
		mediaID := newMedia.ID // Get the ID of the uploaded media

		// Fetch the media record from the database
		var media Models.Media
		if err := c.DB.First(&media, mediaID).Error; err != nil {
			utils.HandleError(w, errors.New("failed to associate media with the post"), http.StatusInternalServerError)
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
			utils.HandleError(w, err, http.StatusInternalServerError)
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
					utils.HandleError(w, err, http.StatusInternalServerError)
					return
				}
			} else {
				utils.HandleError(w, err, http.StatusInternalServerError)
				return
			}
		}

		// Append the tag to the slice
		tags = append(tags, tag)
	}

	// Associate the tags with the post
	if err := c.DB.Model(&newPost).Association("Hashtags").Replace(tags); err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Send email or notification to subscribed users
	// if err := sendPasswordResetEmail(passwordReset.Email, token); err != nil {
	// 	utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	// 	return
	// }
	utils.JSONResponse(w, http.StatusCreated, &utils.APIResponse{
		Type:    "success",
		Message: "The post was created successfully",
	})
}

// Update handles the PUT /posts/{id} request to update a specific post.
//
// This endpoint allows users to update the details of a specific post identified by its ID.
//
// @Summary Update a specific post
// @Description Updates the details of a specific post
// @Tags Posts
// @Param id path int true "Post ID"
// @Accept json
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Router /posts/{id} [put]
func (c PostController) Update(w http.ResponseWriter, r *http.Request) {
	// Handle PUT /postcontroller/{id} request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
	// file, _, err := r.FormFile("file")
	// if err != nil {
	// 	utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	// 	return
	// }
	// defer file.Close()
	// // Generate a unique public ID for the file
	// publicID := uuid.New().String()

	// cloudinaryClient := cloudinary.NewCloudinaryClient()
	// result, err := cloudinaryClient.UploadImage(file, publicID)
	// if err != nil {
	// 	utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	// 	return
	// }
	utils.JSONResponse(w, http.StatusOK, map[string]string{"error": "this is still not implemented"})
}

// UpdateTitle handles the PUT /posts/{id}/title request to update the title of a specific post.
//
// This endpoint allows users to update the title of a specific post identified by its ID.
//
// @Summary Update the title of a specific post
// @Description Updates the title of a specific post
// @Tags Posts
// @Param id path int true "Post ID"
// @Param		post	body		requests.UpdatePostTitleRequest	true	"Post title data"
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Router /posts/{id}/title [put]
func (c PostController) UpdateTitle(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	postID, err := utils.GetParam(r, "id")

	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
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
			utils.HandleError(w, err, http.StatusInternalServerError)
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
			utils.HandleError(w, errors.New("post not found"), http.StatusNotFound)
		} else {
			utils.HandleError(w, result.Error, http.StatusInternalServerError)
		}

		return
	}

	if post.UserID != uint(userID.(float64)) {
		utils.HandleError(w, errors.New("you are not authorized to update post"), http.StatusUnauthorized)
		return
	}
	// Update the post title
	post.Title = updateReq.Title

	result = c.DB.Save(&post)
	if result.Error != nil {
		utils.HandleError(w, result.Error, http.StatusInternalServerError)
		return
	}
	// Return success response
	utils.JSONResponse(w, http.StatusOK,
		&utils.APIResponse{
			Type:    "success",
			Message: "the post title was updated successfully",
			Data:    post,
		})
}

// UpdateBody handles the PUT /posts/{id}/body request to update the body of a specific post.
//
// This endpoint allows users to update the body of a specific post identified by its ID.
//
// @Summary Update the body of a specific post
// @Description Updates the body of a specific post
// @Tags Posts
// @Param id path int true "Post ID"
// @Param post body requests.UpdatePostBodyRequest true "Post body data"
// @Accept json
// @Produce json
// @Success 200 {object} utils.SwaggerSuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Router /posts/{id}/body [put]
func (c PostController) UpdateBody(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	postID, err := utils.GetParam(r, "id")

	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
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
			utils.HandleError(w, err, http.StatusInternalServerError)
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
			utils.HandleError(w, errors.New("post not found"), http.StatusNotFound)
		} else {
			utils.HandleError(w, result.Error, http.StatusInternalServerError)
		}

		return
	}

	if post.UserID != uint(userID.(float64)) {
		utils.HandleError(w, errors.New("you are not authorized to update post"), http.StatusUnauthorized)
		return
	}
	// Update the post title
	post.Body = updateReq.Body

	result = c.DB.Save(&post)
	if result.Error != nil {
		utils.HandleError(w, result.Error, http.StatusInternalServerError)
		return
	}

	// Perform the edit mentions operation
	err = services.EditMentions(c.DB, post.ID, "posts", updateReq.Mentions)
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	// Return success response
	utils.JSONResponse(w, http.StatusOK,
		&utils.APIResponse{
			Type:    "success",
			Message: "the post body was updated successfully",
			Data:    post,
		})
}

// UpdateMedia handles the PUT /posts/{id}/medias request to update the medias of a specific post.
//
// This endpoint allows users to update the medias of a specific post identified by its ID.
//
// @Summary Update the medias of a specific post
// @Description Updates the medias of a specific post
// @Tags Posts
// @Param id path int true "Post ID"
// @Param post body requests.UpdatePostMediaRequest true "Post media data"
// @Accept json
// @Produce json
// @Success 200 {object} utils.SwaggerSuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Router /posts/{id}/medias [put]
func (c PostController) UpdateMedia(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	postID, err := utils.GetParam(r, "id")

	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
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
			utils.HandleError(w, err, http.StatusInternalServerError)
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
			utils.HandleError(w, errors.New("post not found"), http.StatusNotFound)
		} else {
			utils.HandleError(w, result.Error, http.StatusInternalServerError)
		}

		return
	}

	if post.UserID != uint(userID.(float64)) {
		utils.HandleError(w, errors.New("you are not authorized to upate post"), http.StatusUnauthorized)
		return
	}

	// Perform the edit mentions operation
	err = services.EditMedia(c.DB, post.ID, "posts", updateReq.Medias)
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	// Return success response
	utils.JSONResponse(w, http.StatusOK, &utils.APIResponse{
		Type:    "success",
		Message: "post medias was updated successfully",
	})
}

// UpdateHashtag handles the PUT /posts/{id}/hashtags request to update the hashtags of a specific post.
//
// This endpoint allows users to update the hashtags of a specific post identified by its ID.
//
// @Summary Update the hashtags of a specific post
// @Description Updates the hashtags of a specific post
// @Tags Posts
// @Param id path int true "Post ID"
// @Param post body requests.UpdatePostHashtagRequest true "Post hashtag data"
// @Accept json
// @Produce json
// @Success 200 {object} utils.SwaggerSuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Router /posts/{id}/hashtags [put]
func (c PostController) UpdateHashtag(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	postID, err := utils.GetParam(r, "id")

	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
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
			utils.HandleError(w, err, http.StatusInternalServerError)
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
			utils.HandleError(w, err, http.StatusNotFound)
		} else {
			log.Println(err.Error())
			utils.HandleError(w, errors.New("failed to update post tags"), http.StatusInternalServerError)
		}
		return
	}

	// Return success response
	utils.JSONResponse(w, http.StatusOK, &utils.APIResponse{
		Type:    "success",
		Message: "post medias was updated successfully",
	})
}

// UpdatePostSettings handles the PUT /posts/{id}/settings request to update the settings of a specific post.
//
// This endpoint allows users to update the settings of a specific post identified by its ID.
//
// @Summary Update the settings of a specific post
// @Description Updates the settings of a specific post
// @Tags Posts
// @Param id path int true "Post ID"
// @Param post body requests.UpdatePostSettingsRequest true "Post settings data"
// @Accept json
// @Produce json
// @Success 200 {object} utils.SwaggerSuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Router /posts/{id}/settings [put]
func (c *PostController) UpdatePostSettings(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	postID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Parse update data from request body
	var updateReq requests.UpdatePostSettingsRequest
	if err := utils.DecodeJSONBody(w, r, &updateReq); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			utils.JSONResponse(w, mr.Status(), map[string]string{"error": mr.Error()})
		} else {
			log.Print(err.Error())
			utils.HandleError(w, err, http.StatusInternalServerError)
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
			utils.HandleError(w, errors.New("post not found"), http.StatusNotFound)
		} else {
			utils.HandleError(w, result.Error, http.StatusInternalServerError)
		}
		return
	}
	log.Println(updateReq.IsFeatured, updateReq.IsPromoted, updateReq.FeaturedExpiry, updateReq.PromotionExpiry)
	// Update the fields based on the provided update request if the values are not empty or null
	if updateReq.Visibility != "" {
		post.Visibility = updateReq.Visibility
	}
	post.IsPromoted = *updateReq.IsPromoted
	post.IsFeatured = *updateReq.IsFeatured
	post.PromotionExpiry = updateReq.PromotionExpiry
	post.FeaturedExpiry = *updateReq.FeaturedExpiry

	// Save the updated post in the database
	result = c.DB.Save(&post)
	if result.Error != nil {
		utils.HandleError(w, result.Error, http.StatusInternalServerError)
		return
	}

	// Return success response
	utils.JSONResponse(w, http.StatusOK, &utils.APIResponse{
		Type:    "success",
		Message: "post settings was updated successfully",
		Data:    post,
	})
}

// Delete handles the DELETE /posts/{id} request to delete a specific post.
//
// This endpoint allows users to delete a specific post identified by its ID.
//
// @Summary Delete a specific post
// @Description Deletes a specific post
// @Tags Posts
// @Param id path int true "Post ID"
// @Accept json
// @Produce json
// @Success 200 {object} utils.SwaggerSuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Router /posts/{id} [delete]
func (c *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from request parameters
	postID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Get the authenticated user ID from the context
	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Retrieve the post from the database
	var post Models.Post
	result := c.DB.First(&post, postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.HandleError(w, errors.New("post not found"), http.StatusNotFound)
		} else {
			utils.HandleError(w, result.Error, http.StatusInternalServerError)
		}
		return
	}

	// Check if the authenticated user is the owner of the post
	if post.UserID != uint(userID.(float64)) {
		utils.HandleError(w, errors.New("you are not authorized to delete this post"), http.StatusUnauthorized)
		return
	}

	// Delete the post from the database
	if err := c.DB.Delete(&post).Error; err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, &utils.APIResponse{
		Type:    "success",
		Message: "post was deleted successfully",
	})
}
