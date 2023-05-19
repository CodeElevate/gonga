package controllers

import (
	"gonga/app/Models"
	services "gonga/app/Services"
	"gonga/utils"
	"net/http"

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
	// Handle GET /commentcontroller/{id} request

}

func (c CommentController) Create(w http.ResponseWriter, r *http.Request) {
	// Handle POST /commentcontroller request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
}

func (c CommentController) Update(w http.ResponseWriter, r *http.Request) {
	// Handle PUT /commentcontroller/{id} request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
}

func (c CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	// Handle DELETE /commentcontroller/{id} request
}
