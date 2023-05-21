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

func (c LikeController) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	// Validate like request
	if err := utils.ValidateRequest(w, &createReq); err != nil {
		return
	}
	var count int64
	if c.DB.Table(createReq.LikeableType).Where("id = ?", createReq.LikeableID).Count(&count); count == 0 {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "likeabled id doesn't exist."})
		return
	}

	// Check if the user has already liked the record
	var existingLike Models.Like
	if err := c.DB.Where("user_id = ? AND likeable_id = ? AND likeable_type = ?",
		userID, createReq.LikeableID, createReq.LikeableType).First(&existingLike).Error; err == nil {
		// User has already liked the record, perform unlike operation
		if err := c.DB.Unscoped().Delete(&existingLike).Error; err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "unliked successfully!"})
		return
	}

	like := &Models.Like{
		UserID:       uint(userID.(float64)),
		LikeableID:   createReq.LikeableID,
		LikeableType: createReq.LikeableType,
	}
	// Save the like to the database
	if err := c.DB.Create(&like).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "Failed to save like", "error": err.Error()})
		return
	}
	// Return success response
	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "like was created successfully!"})
}

func (c LikeController) Update(w http.ResponseWriter, r *http.Request) {
	// Handle PUT /likecontroller/{id} request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
}

func (c LikeController) Delete(w http.ResponseWriter, r *http.Request) {
	likeID, err := utils.GetParam(r, "id") // Get the like ID from the URL parameter
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	userID, err := utils.GetUserIDFromContext(r.Context()) // Get the user ID from the context

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Check if the like exists
	var like Models.Like
	if err := c.DB.First(&like, likeID).Error; err != nil {
		utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "Like not found"})
		return
	}

	// Check if the current user is the owner of the like
	if like.UserID != uint(userID.(float64)) {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Delete the like
	if err := c.DB.Delete(&like).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Like deleted successfully"})
}
