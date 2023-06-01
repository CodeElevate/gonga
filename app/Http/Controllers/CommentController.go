package controllers

import (
	"errors"
	requests "gonga/app/Http/Requests"
	"gonga/app/Models"
	services "gonga/app/Services"

	// services "gonga/app/Services"
	"gonga/utils"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type CommentController struct {
	DB *gorm.DB
}

// Index handles the GET /posts/{id}/comments request to retrieve comments for a post.
//
// This endpoint allows users to retrieve comments for a specific post based on its ID.
//
//	@Summary		Get comments for a post
//	@Description	Retrieves comments for a specific post
//	@Tags			Comments
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	utils.SwaggerPagination
//	@Failure		400	{object}	utils.SwaggerErrorResponse
//	@Failure		500	{object}	utils.SwaggerErrorResponse
//	@Router			/posts/{id}/comments [get]
func (c CommentController) Index(w http.ResponseWriter, r *http.Request) {
	// Handle GET /postcontroller/{id} request
	postID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Handle GET /commentcontroller request
	var comments []Models.Comment
	var response utils.APIResponse

	// Apply the where condition to filter comments by postID and parentID
	db := c.DB.Where("post_id = ? AND parent_id IS NULL", postID)

	// Apply pagination and retrieve paginated comments
	paginationScope, err := utils.Paginate(r, db, &comments, &response, "User", "Mentions.User", "Childrens", "Likes")
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Apply the pagination scope to the filtered query
	db = paginationScope(db)

	// Retrieve the paginated comments
	if err := db.Find(&comments).Error; err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Set the items value in the pagination struct
	response.Data = comments
	response.Type = "success"
	response.Message = "data fetched successfully!"

	// Send the response with the pagination struct
	utils.JSONResponse(w, http.StatusOK, response)
}

// Show handles the GET /comments/{id} request to retrieve a specific comment.
//
// This endpoint retrieves a specific comment by its ID.
//
//	@Summary		Get a comment
//	@Description	Retrieves a specific comment by its ID
//	@Tags			Comments
//	@Param			id	path	string	true	"Comment ID"
//	@Produce		json
//	@Success		200	{object}	utils.SwaggerSuccessResponse
//	@Failure		400	{object}	utils.SwaggerErrorResponse
//	@Failure		404	{object}	utils.SwaggerErrorResponse
//	@Failure		500	{object}	utils.SwaggerErrorResponse
//	@Router			/comments/{id} [get]
func (c CommentController) Show(w http.ResponseWriter, r *http.Request) {
	commentID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Create a variable to hold the comment
	var comment Models.Comment

	// Retrieve the comment with the specified ID from the database
	if err := c.DB.Preload("Childrens").First(&comment, commentID).Error; err != nil {
		// If the comment is not found, return a not found response
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(w, err, http.StatusNotFound, "comment not found")
			return
		}

		// If an error occurs during the database query, return an internal server error response
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Send the comment as a response
	utils.JSONResponse(w, http.StatusOK, utils.APIResponse{
		Data:    comment,
		Type:    "success",
		Message: "Comment retrieved successfully!",
	})
}

// Create handles the POST /posts/{id}/comments request to create a new comment for a post.
//
// This endpoint allows authenticated users to create a new comment for a specific post.
//
//	@Summary		Create a new comment
//	@Description	Creates a new comment for a specific post
//	@Tags			Comments
//	@Param			id				path	string	true	"Post ID"
//	@Param			Authorization	header	string	true	"Bearer token"
//	@Accept			json
//	@Produce		json
//	@Param			body	body		requests.CreateCommentRequest	true	"Comment data"
//	@Success		200		{object}	utils.SwaggerSuccessResponse
//	@Failure		400		{object}	utils.SwaggerErrorResponse
//	@Failure		401		{object}	utils.SwaggerErrorResponse
//	@Failure		500		{object}	utils.SwaggerErrorResponse
//	@Router			/posts/{id}/comments [post]
func (c CommentController) Create(w http.ResponseWriter, r *http.Request) {
	postIDStr, err := utils.GetParam(r, "id")
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Parse update request from request body
	var createReq requests.CreateCommentRequest

	if err := utils.DecodeJSONBody(w, r, &createReq); err != nil {
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
	if err := utils.ValidateRequest(w, &createReq); err != nil {
		return
	}

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Convert postID from string to uint
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		utils.HandleError(w, errors.New("invalid post ID"), http.StatusBadRequest)
		return
	}
	// Check if the parent comment exists
	if createReq.ParentID != nil {
		var parentComment Models.Comment
		if err := c.DB.First(&parentComment, createReq.ParentID).Error; err != nil {
			utils.HandleError(w, errors.New("invalid parent comment ID"), http.StatusBadRequest)
			return
		}
	}
	// Create a new Comment instance
	newComment := Models.Comment{
		UserID:   uint(userID.(float64)),
		PostID:   uint(postID),
		Body:     createReq.Body,
		ParentID: createReq.ParentID,
	}

	// Insert the comment into the database
	if err := c.DB.Create(&newComment).Error; err != nil {
		log.Println(err.Error())
		utils.HandleError(w, errors.New("failed to create comment"), http.StatusInternalServerError)
		return
	}

	// Create mentions for the comment
	for _, mentionedUser := range createReq.Mentions {
		log.Println(mentionedUser.UserID)
		mention := &Models.Mention{
			UserID:    mentionedUser.UserID,
			OwnerID:   newComment.ID,
			OwnerType: "comments",
		}
		// Save the mention to the database
		if err := c.DB.Create(&mention).Error; err != nil {
			utils.HandleError(w, err, http.StatusInternalServerError)
			return
		}
	}

	// Send the created comment as a response
	utils.JSONResponse(w, http.StatusOK, &utils.APIResponse{
		Type:    "success",
		Message: "comment created successfully!",
	})
}

// Update handles the PUT /comments/{id} request to update a specific comment.
//
// This endpoint updates a specific comment by its ID.
//
//	@Summary		Update a comment
//	@Description	Updates a specific comment by its ID
//	@Tags			Comments
//	@Param			id	path	string	true	"Comment ID"
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer token"
//	@Param			body			body		requests.UpdateCommentRequest	true	"Update Comment Request"
//	@Success		200				{object}	utils.SwaggerSuccessResponse
//	@Failure		400				{object}	utils.SwaggerErrorResponse
//	@Failure		401				{object}	utils.SwaggerErrorResponse
//	@Failure		404				{object}	utils.SwaggerErrorResponse
//	@Failure		500				{object}	utils.SwaggerErrorResponse
//	@Router			/comments/{id} [put]
func (c CommentController) Update(w http.ResponseWriter, r *http.Request) {
	// Extract the comment ID from the URL path parameters
	commentID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Parse the request body into the UpdateCommentRequest struct
	var updateReq requests.UpdateCommentRequest

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

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Check if the comment exists and belongs to the authenticated user
	var comment Models.Comment
	if err := c.DB.Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(w, errors.New("comment not found"), http.StatusNotFound)
		} else {
			utils.HandleError(w, err, http.StatusInternalServerError, "failed to fetch comment")
		}
		return
	}

	// Update the comment body
	comment.Body = updateReq.Body

	// Save the updated comment
	if err := c.DB.Save(&comment).Error; err != nil {
		utils.HandleError(w, errors.New("failed to update comment"), http.StatusInternalServerError)
		return
	}

	// Perform the edit mentions operation
	err = services.EditMentions(c.DB, comment.ID, "comments", updateReq.Mentions)
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	// Send the updated comment in the response
	utils.JSONResponse(w, http.StatusOK, &utils.APIResponse{
		Type:    "success",
		Message: "comment updated successfully!",
		Data:    comment,
	})

}

// Delete handles the DELETE /comments/{id} request to delete a specific comment.
//
// This endpoint deletes a specific comment by its ID.
//
//	@Summary		Delete a comment
//	@Description	Deletes a specific comment by its ID
//	@Tags			Comments
//	@Param			id	path	string	true	"Comment ID"
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer token"
//	@Success		200				{object}	utils.SwaggerSuccessResponse
//	@Failure		400				{object}	utils.SwaggerErrorResponse
//	@Failure		401				{object}	utils.SwaggerErrorResponse
//	@Failure		404				{object}	utils.SwaggerErrorResponse
//	@Failure		500				{object}	utils.SwaggerErrorResponse
//	@Router			/comments/{id} [delete]
func (c CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract the comment ID from the URL path parameters
	commentID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Check if the comment exists
	var comment Models.Comment
	if err := c.DB.First(&comment, commentID).Error; err != nil {
		utils.HandleError(w, errors.New("comment not found"), http.StatusNotFound)
		return
	}

	// Delete the comment
	if err := c.DB.Delete(&comment).Error; err != nil {
		utils.HandleError(w, errors.New("failed to delete comment"), http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, utils.APIResponse{
		Type: "success",
		Message: "comment deleted successfully",
	})
}
