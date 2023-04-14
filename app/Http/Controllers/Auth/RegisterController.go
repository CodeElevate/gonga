package auth

import (
	"gonga/app/Models"
	"gonga/utils"
	"net/http"

	"gorm.io/gorm"
)

type RegisterController struct {
	DB *gorm.DB
}

type RegisterUser struct {
	Username string `json:"username" validate:"required,min=3,max=5"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type RegisterResponse struct {
	Token   string `json:"token"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

func (c RegisterController) Index(w http.ResponseWriter, r *http.Request) {
	// Handle GET /registercontroller request
}

func (c RegisterController) Show(w http.ResponseWriter, r *http.Request) {
	// Handle GET /registercontroller/{id} request
}

func (c RegisterController) Create(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user RegisterUser
	if err := utils.DecodeRequestBody(r, &user); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// Validate user data
	if err := utils.ValidateRequest(w, &user); err != nil {
		return
	}

	// Create user in database
	hashedPassword, err := utils.HashPassword(user.Password)
	
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	newUser := Models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}
	result := c.DB.Create(&newUser)
	if result.Error != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(int(newUser.ID))
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Send response
	response := RegisterResponse{
		Token:   token,
		UserID:  int(newUser.ID),
		Message: "Registration successful",
	}
	utils.JSONResponse(w, http.StatusOK, response)
}

func (c RegisterController) Update(w http.ResponseWriter, r *http.Request) {
	// Handle PUT /registercontroller/{id} request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
}

func (c RegisterController) Delete(w http.ResponseWriter, r *http.Request) {
	// Handle DELETE /registercontroller/{id} request
}
