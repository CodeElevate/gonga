package controllers

import (
	"gonga/utils"
	"net/http"

	requests "gonga/app/Http/Requests"
	"gonga/app/Models"

	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

//	@Summary		Get a list of users with pagination
//	@Description	Retrieve a paginated list of users
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int					false	"Page number for pagination"
//	@Param			per_page	query		int					false	"Number of items per page"
//	@Success		200			{object}	utils.APIResponse	"success"
//	@Failure		500			{object}	utils.APIResponse	"error"
//	@Router			/users [GET]
func (uc UserController) Index(w http.ResponseWriter, r *http.Request) {
	var users []Models.User
	var response utils.APIResponse

	paginationScope, err := utils.Paginate(r, uc.DB, &users, &response)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	db := paginationScope(uc.DB)
	if err := db.Find(&users).Error; err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Set the items value in the pagination struct
	response.Data = users
	response.Type = "success"
	response.Message = "data retrieved successfully"

	// Send the response with the pagination struct
	utils.JSONResponse(w, http.StatusOK, response)

}

//	@Summary		Get a user by username
//	@Description	Retrieve a user by their username
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string				true	"Username of the user to retrieve"
//	@Success		200			{object}	utils.APIResponse	"success"
//	@Failure		400			{object}	utils.APIResponse	"error"
//	@Failure		404			{object}	utils.APIResponse	"error"
//	@Router			/users/{username} [GET]
func (uc UserController) Show(w http.ResponseWriter, r *http.Request) {
	username, err := utils.GetParam(r, "username")
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// Fetch user from the database
	var user Models.User
	if err := uc.DB.Where("username = ?", username).Preload("FollowersList").Preload("FollowingList").First(&user).Error; err != nil {
		// User not found, return error response
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Return successful response with the user data
	utils.JSONResponse(w, http.StatusOK, utils.APIResponse{
		Type: "success",
		Data: user,
	})
}

func (uc UserController) Create(w http.ResponseWriter, r *http.Request) {
	// Handle POST /users request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
}

//	@Summary		Update a user
//	@Description	Update the details of a user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string						true	"Username of the user to update"
//	@Param			updateReq	body		requests.UpdateUserRequest	true	"Update request body"
//	@Success		200			{object}	utils.APIResponse			"success"
//	@Failure		400			{object}	utils.APIResponse			"error"
//	@Failure		404			{object}	utils.APIResponse			"error"
//	@Failure		500			{object}	utils.APIResponse			"error"
//	@Router			/users/{username} [PUT]
//	@Security		BearerAuth
func (uc UserController) Update(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request path
	username, err := utils.GetParam(r, "username")
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	var user Models.User
	if err := uc.DB.Where("username = ?", username).First(&user).Error; err != nil {
		utils.HandleError(w, err, http.StatusNotFound)
		return
	}

	// Parse update request from request body
	var updateReq requests.UpdateUserRequest
	if err := utils.DecodeRequestBody(r, &updateReq); err != nil {
		utils.HandleError(w, err, http.StatusBadRequest)
		return
	}

	// Validate the update request
	if err := utils.ValidateRequest(w, &updateReq); err != nil {
		return
	}

	// Only update fields that are present in the update request
	if updateReq.FirstName != "" {
		user.FirstName = updateReq.FirstName
	}
	if updateReq.LastName != "" {
		user.LastName = updateReq.LastName
	}
	if updateReq.AvatarURL != "" {
		user.AvatarURL = updateReq.AvatarURL
	}
	if updateReq.Bio != "" {
		user.Bio = updateReq.Bio
	}
	if updateReq.Gender != "" {
		user.Gender = updateReq.Gender
	}
	if updateReq.MobileNo != "" {
		user.MobileNo = updateReq.MobileNo
	}
	if updateReq.MobileNoCode != "" {
		user.MobileNoCode = updateReq.MobileNoCode
	}
	if !updateReq.Birthday.IsZero() {
		user.Birthday = updateReq.Birthday
	}
	if updateReq.Country != "" {
		user.Country = updateReq.Country
	}
	if updateReq.City != "" {
		user.City = updateReq.City
	}
	if updateReq.BackgroundImageURL != "" {
		user.BackgroundImageURL = updateReq.BackgroundImageURL
	}
	if updateReq.WebsiteURL != "" {
		user.WebsiteURL = updateReq.WebsiteURL
	}
	if updateReq.Occupation != "" {
		user.Occupation = updateReq.Occupation
	}
	if updateReq.Education != "" {
		user.Education = updateReq.Education
	}

	// Save updated user to the database
	if err := uc.DB.Save(&user).Error; err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Send success response with updated user information
	utils.JSONResponse(w, http.StatusOK, utils.APIResponse{
		Type: "success",
		Data: user,
	})
}

//	@Summary		Delete a user
//	@Description	Delete a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID of the user to delete"
//	@Success		204	"success"
//	@Failure		404	{object}	utils.APIResponse	"error"
//	@Failure		500	{object}	utils.APIResponse	"error"
//	@Router			/users/{id} [DELETE]
//	@Security		BearerAuth
func (uc UserController) Delete(w http.ResponseWriter, r *http.Request) {
	// Handle DELETE /users/{id} request
}
