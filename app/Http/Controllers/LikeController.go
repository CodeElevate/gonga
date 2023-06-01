package controllers

import (
	"errors"
	requests "gonga/app/Http/Requests"
	"gonga/app/Models"
	"gonga/utils"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type LikeController struct {
	DB *gorm.DB
}

func (c LikeController) Index(w http.ResponseWriter, r *http.Request) {
	// Handle GET /likecontroller request
}

func (c LikeController) Show(w http.ResponseWriter, r *http.Request) {
	// Handle GET /likecontroller/{id} request
}

// Create handles the POST /likes request to create a new like.
//
// This endpoint allows users to create a new like for a specific likeable item.
//
//	@Summary		Create a new like
//	@Description	Creates a new like
//	@Tags			Likes
//	@Accept			json
//	@Produce		json
//	@Param			like	body		requests.CreateLikeRequest	true	"Like data"
//	@Success		200		{object}	utils.SwaggerSuccessResponse
//	@Failure		400		{object}	utils.SwaggerErrorResponse
//	@Failure		401		{object}	utils.SwaggerErrorResponse
//	@Failure		500		{object}	utils.SwaggerErrorResponse
//	@Router			/likes [post]
func (c LikeController) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	// Parse like data from request body
	var createReq requests.CreateLikeRequest
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
	// Validate like request
	if err := utils.ValidateRequest(w, &createReq); err != nil {
		return
	}
	var count int64
	if c.DB.Table(createReq.LikeableType).Where("id = ?", createReq.LikeableID).Count(&count); count == 0 {
		utils.HandleError(w, errors.New("likeabled id doesn't exist"), http.StatusNotFound)
		return
	}

	// Check if the user has already liked the record
	var existingLike Models.Like
	if err := c.DB.Where("user_id = ? AND likeable_id = ? AND likeable_type = ?",
		userID, createReq.LikeableID, createReq.LikeableType).First(&existingLike).Error; err == nil {
		// User has already liked the record, perform unlike operation
		if err := c.DB.Unscoped().Delete(&existingLike).Error; err != nil {
			utils.HandleError(w, err, http.StatusInternalServerError)
			return
		}
		utils.JSONResponse(w, http.StatusOK, utils.APIResponse{
			Type:    "success",
			Message: "unliked successfully!",
		})
		return
	}

	like := &Models.Like{
		UserID:       uint(userID.(float64)),
		LikeableID:   createReq.LikeableID,
		LikeableType: createReq.LikeableType,
	}
	// Save the like to the database
	if err := c.DB.Create(&like).Error; err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "failed to save like")
		return
	}
	// Return success response
	utils.JSONResponse(w, http.StatusOK, utils.APIResponse{
		Type:    "success",
		Message: "like was created successfully!",
	})
}

func (c LikeController) Update(w http.ResponseWriter, r *http.Request) {
	// Handle PUT /likecontroller/{id} request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
}

// Delete handles the DELETE /likes/{id} request to delete a like.
//
// This endpoint allows users to delete a like based on its ID.
//
//	@Summary		Delete a like
//	@Description	Deletes a like
//	@Tags			Likes
//	@Param			id	path		string	true	"Like ID"
//	@Success		200	{object}	utils.SwaggerSuccessResponse
//	@Failure		400	{object}	utils.SwaggerErrorResponse
//	@Failure		401	{object}	utils.SwaggerErrorResponse
//	@Failure		404	{object}	utils.SwaggerErrorResponse
//	@Failure		500	{object}	utils.SwaggerErrorResponse
//	@Router			/likes/{id} [delete]
func (c LikeController) Delete(w http.ResponseWriter, r *http.Request) {
	likeID, err := utils.GetParam(r, "id") // Get the like ID from the URL parameter
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}
	userID, err := utils.GetUserIDFromContext(r.Context()) // Get the user ID from the context

	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Check if the like exists
	var like Models.Like
	if err := c.DB.First(&like, likeID).Error; err != nil {
		utils.HandleError(w, errors.New("like not found"), http.StatusNotFound)
		return
	}

	// Check if the current user is the owner of the like
	if like.UserID != uint(userID.(float64)) {
		utils.HandleError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	// Delete the like
	if err := c.DB.Delete(&like).Error; err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Return success response
	utils.JSONResponse(w, http.StatusOK, utils.APIResponse{
		Type:    "success",
		Message: "like deleted successfully",
	})
}
