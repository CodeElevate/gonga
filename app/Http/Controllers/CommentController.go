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

func (c CommentController) Index(w http.ResponseWriter, r *http.Request) {
	// Handle GET /postcontroller/{id} request
	postID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Apply the pagination scope to the filtered query
	db = paginationScope(db)

	// Retrieve the paginated comments
	if err := db.Find(&comments).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Set the items value in the pagination struct
	response.Data = comments
	response.Type = "success"
	response.Message = "data fetched successfully!"

	// Send the response with the pagination struct
	utils.JSONResponse(w, http.StatusOK, response)
}


func (c CommentController) Show(w http.ResponseWriter, r *http.Request) {
	commentID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Create a variable to hold the comment
	var comment Models.Comment

	// Retrieve the comment with the specified ID from the database
	if err := c.DB.Preload("Childrens").First(&comment, commentID).Error; err != nil {
		// If the comment is not found, return a not found response
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Comment not found"})
			return
		}

		// If an error occurs during the database query, return an internal server error response
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the comment as a response
	utils.JSONResponse(w, http.StatusOK, utils.APIResponse{
		Data: comment,
		Type: "success",
	})
}

// AddComment adds a comment to a post.
//	@Summary	Add a comment to a post
//	@Tags		Comments
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int								true	"Post ID"
//	@Param		comment	body		requests.CreateCommentRequest	true	"Comment data"
//	@Success	200		{object}	map[string]string
//	@Failure	400		{object}	map[string]string
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/posts/{id}/comments [post]
func (c CommentController) Create(w http.ResponseWriter, r *http.Request) {
	postIDStr, err := utils.GetParam(r, "id")
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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

	// Convert postID from string to uint
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}
	// Check if the parent comment exists
	if createReq.ParentID != nil {
		var parentComment Models.Comment
		if err := c.DB.First(&parentComment, createReq.ParentID).Error; err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid parent comment ID"})
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
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create comment"})
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
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	// Send the created comment as a response
	utils.JSONResponse(w, http.StatusOK, &utils.ControllerResponse{
		Message: "comment created successfully!",
		Data:    newComment,
	})
}

// UpdateComment updates a comment on a post.
//
//	@Summary	Update a comment on a post
//	@Tags		Comments
//	@Accept		json
//	@Produce	json
//	@Param		id			path		int								true	"Post ID"
//	@Param		comment_id	path		int								true	"Comment ID"
//	@Param		comment		body		requests.UpdateCommentRequest	true	"Comment data"
//	@Success	200			{object}	map[string]string
//	@Failure	400			{object}	map[string]string
//	@Failure	404			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/posts/{id}/comments/{comment_id} [put]
func (c CommentController) Update(w http.ResponseWriter, r *http.Request) {
	// Extract the comment ID from the URL path parameters
	commentID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	// Validate post request
	if err := utils.ValidateRequest(w, &updateReq); err != nil {
		return
	}

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Check if the comment exists and belongs to the authenticated user
	var comment Models.Comment
	if err := c.DB.Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Comment not found"})
		} else {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch comment"})
		}
		return
	}

	// Update the comment body
	comment.Body = updateReq.Body

	// Save the updated comment
	if err := c.DB.Save(&comment).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update comment"})
		return
	}

	// Perform the edit mentions operation
	err = services.EditMentions(c.DB, comment.ID, "comments", updateReq.Mentions)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Send the updated comment in the response
	utils.JSONResponse(w, http.StatusOK, &utils.ControllerResponse{
		Message: "Comment updated successfully!",
		Data:    comment,
	})

}

// DeleteComment deletes a comment from a post.
//
//	@Summary	Delete a comment from a post
//	@Tags		Comments
//	@Produce	json
//	@Param		id			path		int	true	"Post ID"
//	@Param		comment_id	path		int	true	"Comment ID"
//	@Success	200			{object}	map[string]string
//	@Failure	404			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/posts/{id}/comments/{comment_id} [delete]
func (c CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract the comment ID from the URL path parameters
	commentID, err := utils.GetParam(r, "id")
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Check if the comment exists
	var comment Models.Comment
	if err := c.DB.First(&comment, commentID).Error; err != nil {
		utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Comment not found"})
		return
	}

	// Delete the comment
	if err := c.DB.Delete(&comment).Error; err != nil {

		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete comment"})
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"error": "Comment deleted successfully"})
}
