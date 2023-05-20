package controllers

import (
	"errors"
	requests "gonga/app/Http/Requests"
	"gonga/app/Models"
	services "gonga/app/Services"
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
	var pagination utils.Pagination

	paginationScope, err := utils.Paginate(r, c.DB, &comments, &pagination, "User", "Mentions.User")
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	db := paginationScope(c.DB)

	if err := db.Where("post_id = ? AND parent_id IS NULL", postID).Find(&comments).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	for i := range comments {
		services.LoadNestedComments(&comments[i], c.DB)
	}
	// Set the items value in the pagination struct
	pagination.Items = comments

	// Send the response with the pagination struct
	utils.JSONResponse(w, http.StatusOK, pagination)
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
	if err := c.DB.First(&comment, commentID).Error; err != nil {
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
	utils.JSONResponse(w, http.StatusOK, comment)
}

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
	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "comment created successfully!"})
}

func (c CommentController) Update(w http.ResponseWriter, r *http.Request) {
	// Handle PUT /commentcontroller/{id} request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
}

func (c CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	// Handle DELETE /commentcontroller/{id} request
}
