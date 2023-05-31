package auth

import (
	"errors"
	requests "gonga/app/Http/Requests/Auth"
	responses "gonga/app/Http/Responses/Auth"
	"gonga/utils"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type LoginController struct {
	DB *gorm.DB
}

func (c LoginController) Index(w http.ResponseWriter, r *http.Request) {
	// Handle GET /logincontroller request
}

func (c LoginController) Show(w http.ResponseWriter, r *http.Request) {
	// Handle GET /logincontroller/{id} request
}

// Create handles the POST /login request for user login.
//
// This endpoint allows users to log in by providing their username and password.
//
// @Summary User login
// @Description Logs in a user with the provided credentials
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginRequest body requests.LoginRequest true "Login credentials"
// @Success 200 {object} responses.LoginResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Router /login [post]
func (c LoginController) Create(w http.ResponseWriter, r *http.Request) {
	// Handle POST /logincontroller request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
	// Parse request body
	var user requests.LoginRequest

	if err := utils.DecodeJSONBody(w, r, &user); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			utils.JSONResponse(w, mr.Status(), map[string]string{"error": mr.Error()})
		} else {
			log.Print(err.Error())
			utils.HandleError(w, err, http.StatusInternalServerError)
		}
		return
	}
	// Validate user data

	if err := utils.ValidateRequest(w, &user); err != nil {
		return
	}
	// Check credentials and get user ID from database
	userID, err := utils.Authenticate(user.Username, user.Password, c.DB)
	if err != nil {
		utils.HandleError(w, errors.New("invalid username or password"), http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(userID)
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	// Send response
	response := responses.LoginResponse{
		Token:   token,
		UserID:  userID,
		Message: "Login successful",
	}
	utils.JSONResponse(w, http.StatusOK, response)
}

func (c LoginController) Update(w http.ResponseWriter, r *http.Request) {
	// Handle PUT /logincontroller/{id} request
	// You can get the request body by reading from r.Body
	// You can send a response by writing to w
}

func (c LoginController) Delete(w http.ResponseWriter, r *http.Request) {
	// Handle DELETE /logincontroller/{id} request
}
