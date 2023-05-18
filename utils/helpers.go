package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Env gets the value of an environment variable.
//
// If the environment variable with the specified key exists, its value is returned.
// Otherwise, the fallback value is returned.
//
// Example usage:
//
//	port := Env("PORT", "8080")
//
// Parameters:
//   - key (string): The name of the environment variable to get the value of.
//   - fallback (string): The fallback value to return if the environment variable does not exist.
//
// Returns:
//   - string: The value of the environment variable or the fallback value.
func Env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// EnvInt gets the integer value of an environment variable.
//
// If the environment variable with the specified key exists and its value is an integer,
// its integer value is returned. Otherwise, the fallback value is returned.
//
// Example usage:
//
//	maxConnections := EnvInt("MAX_CONNECTIONS", 100)
//
// Parameters:
//   - key (string): The name of the environment variable to get the value of.
//   - fallback (int): The fallback value to return if the environment variable does not exist or its value is not an integer.
//
// Returns:
//   - int: The integer value of the environment variable or the fallback value.
func EnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		num, err := strconv.Atoi(value)
		if err != nil {
			// handle error
			fmt.Println("Error:", err)
			os.Exit(0)
		}
		return num
	}
	return fallback
}

// JSONResponse sends a JSON response with the specified status code and data.
// It sets the Content-Type header to "application/json", writes the status code to
// the response writer, and marshals the data to JSON format. If an error occurs
// during marshalling, it returns a 500 Internal Server Error with the error message.
//
// Example usage:
//
//	data := map[string]string{"message": "Hello, world!"}
//	JSONResponse(w, http.StatusOK, data)
//
// Parameters:
//   - w (http.ResponseWriter): The response writer to write the JSON response to.
//   - statusCode (int): The status code to set for the response.
//   - data (interface{}): The data to be marshalled to JSON format and sent as the response.
//
// Returns:
//   - void
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set status code
	w.WriteHeader(statusCode)

	// Marshal data to JSON format
	responseJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response JSON to response writer
	w.Write(responseJSON)
}

// ErrorResponse is a struct that represents the response format for validation errors.
// It contains a map of field names and their corresponding error messages.
//
// Example usage:
//
//	errResponse := ErrorResponse{
//	   Errors: map[string]string{
//	       "username": "Username is already taken",
//	       "email": "Email address is invalid",
//	   },
//	}
//
// Fields:
//   - Errors (map[string]string): The map of field names and their corresponding error messages.
type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

// ValidateRequest validates the given request data against the specified validation rules
// using the Go Playground Validator library, and sends a formatted error response if validation fails.
// If validation is successful, it returns nil.
//
// Example usage:
//
//	type User struct {
//	    Name  string `json:"name" validate:"required"`
//	    Email string `json:"email" validate:"required,email"`
//	}
//
//	func CreateUser(w http.ResponseWriter, r *http.Request) {
//	    var user User
//	    decoder := json.NewDecoder(r.Body)
//	    if err := decoder.Decode(&user); err != nil {
//	        JSONResponse(w, http.StatusBadRequest, ErrorResponse{Errors: map[string]string{"message": "Invalid request payload"}})
//	        return
//	    }
//
//	    // Validate the request payload using the `User` struct as validation rules
//	    if err := ValidateRequest(w, user); err != nil {
//	        return
//	    }
//
//	    // Create the user
//	    // ...
//	}
//
// Parameters:
//   - w (http.ResponseWriter): The response writer to send the error response to.
//   - rules (interface{}): The validation rules to apply to the request data. The validation rules
//     must be in the format supported by the Go Playground Validator library.
//
// Returns:
//   - error: If validation fails, returns an error with the validation errors. Otherwise, returns nil.
func ValidateRequest(w http.ResponseWriter, rules interface{}) error {
	// Validate the form data using the validator library
	validate := validator.New()
	// Send a formatted error response if validation fails
	if err := validate.Struct(rules); err != nil {
		// Extract field names and error messages from the validation errors
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
			errors[e.Field()] = ValidationMessage(e)
		}
		// Create the error response
		response := ErrorResponse{Errors: errors}
		// Send the error response
		JSONResponse(w, http.StatusBadRequest, response)
		
		return err
	}
	return nil
}

// GetPaginationParams extracts pagination parameters from the provided HTTP request's query string.
// If the "page" or "per_page" parameters are not present in the query string, this function will use the provided default values.
// It returns two integers: the page number and the number of items per page.
//
// Example usage:
//
//	page, perPage := GetPaginationParams(r, 1, 10)
func GetPaginationParams(r *http.Request, defaultPage, defaultPerPage int) (page int, perPage int) {
	// Get page and per_page parameters from query string or use default values
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		page = defaultPage
	} else {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			page = defaultPage
		} else {
			page = pageInt
		}
	}

	perPageStr := r.URL.Query().Get("per_page")
	if perPageStr == "" {
		perPage = defaultPerPage
	} else {
		perPageInt, err := strconv.Atoi(perPageStr)
		if err != nil {
			perPage = defaultPerPage
		} else {
			perPage = perPageInt
		}
	}

	return page, perPage
}

func GetParam(r *http.Request, paramName string) (string, error) {
	params := mux.Vars(r)
	val, ok := params[paramName]
	if !ok {
		return "", fmt.Errorf("missing parameter %s", paramName)
	}
	return val, nil
}
