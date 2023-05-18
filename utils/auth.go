package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gonga/app/Models"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/gddo/httputil/header"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

// IsAuthenticate checks if the user is authenticated by verifying the JWT token present in the request header.
//
// It expects the JWT token to be present in the "token" field of the request header. If it's not found, or if the token
// is invalid or expired, this function returns false. Otherwise, it returns true indicating that the user is authenticated.
//
// Example usage:
//
//	isAuthenticated := IsAuthenticate(r)
//	if isAuthenticated {
//	    // perform authenticated operations
//	} else {
//	    // redirect to login page or send unauthorized response
//	}
//
// Parameters:
//
//	r (*http.Request): The HTTP request object containing the token in its header.
//
// Returns:
//
//	bool: A boolean value indicating whether the user is authenticated or not.
func IsAuthenticate(r *http.Request) (bool, interface{}) {
	// Check if user is authenticated
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return false, nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Provide a signing key
		return []byte(Env("APP_KEY", "my-secret-key")), nil
	})
	if err != nil {
		return false, nil
	}
	if !token.Valid {
		return false, nil
	}

	// Extract the user ID from the token's claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	userID, exists := claims["userID"]
	if !exists {
		return false, nil
	}

	return true, userID
}

func ExtractUserIDFromToken(tokenString string) (int, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Add your logic to provide the token's signing key here
		// For example, you can use a shared secret key or fetch the key from a database
		// Make sure to return the correct signing key based on the token's signing method
		// In this example, we assume the token is signed using an HMAC algorithm
		return []byte(Env("APP_KEY", "my-secret-key")), nil
	})

	if err != nil {
		return 0, err
	}

	// Verify the token's validity
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Extract the user ID from the token's claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["userID"].(int)
	if !ok {
		return 0, errors.New("user ID not found in token claims")
	}

	return userID, nil
}

func GetUserIDFromContext(ctx context.Context) (interface{}, error) {
	const (
		userIDKey ContextKey = "userID"
	)
	userID := ctx.Value(userIDKey)

	if userID != nil {
		switch v := userID.(type) {
		case uint:
			return v, nil
		case float64:
			return v, nil
		case string:
			// Convert the string to the appropriate type (e.g., UUID)
			// Handle different type conversions as needed
			userUUID, err := uuid.Parse(v)
			if err != nil {
				return nil, err
			}
			return userUUID, nil
		default:
			return nil, fmt.Errorf("unsupported user ID type: %T", userID)
		}
	}
	return nil, errors.New("user ID not found in context")
}

// Authenticate authenticates a user by checking the provided username and password against the database.
//
// It fetches the user with the given username from the database using the provided *gorm.DB object. If the user is found,
// it verifies the provided password against the user's hashed password using the VerifyPassword() function. If the password
// is correct, it returns the user ID as an integer. Otherwise, it returns an error indicating the reason for the failure.
// Example usage:
//
//	userID, err := Authenticate("john.doe", "my-password", db)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// perform authenticated operations with userID
//
// Parameters:
//
//	username (string): The username of the user to be authenticated.
//	password (string): The password of the user to be authenticated.
//	db (*gorm.DB): The *gorm.DB object to be used for fetching the user from the database.
//
// Returns:
//
//	int: The user ID as an integer, if authentication is successful.
//	error: An error, if any, that occurred during the authentication process.
func Authenticate(username, password string, db *gorm.DB) (int, error) {
	// Get user by username from database
	var user Models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}

	// Compare passwords
	if err := VerifyPassword(user.Password, password); err != nil {
		log.Println(err)
		return 0, err
	}

	return int(user.ID), nil
}

// UserExistsWithEmail checks if a user with the given email address exists in the database.
//
// It searches for the user with the provided email address in the database using the provided *gorm.DB object.
// If a user with the given email exists in the database, it returns true, otherwise it returns false.
//
// Example usage:
//
//	email := "john.doe@example.com"
//	exists := UserExistsWithEmail(db, email)
//	if exists {
//	    fmt.Println("A user with the email", email, "already exists.")
//	} else {
//	    fmt.Println("The email", email, "is available.")
//	}
//
// Parameters:
//
//	db (*gorm.DB): The *gorm.DB object to be used for fetching the user from the database.
//	email (string): The email address to be checked.
//
// Returns:
//
//	bool: True, if a user with the provided email exists in the database, otherwise false.
func UserExistsWithEmail(db *gorm.DB, email string) bool {
	var user Models.User
	result := db.Where("email = ?", email).First(&user)
	return result.Error == nil && result.RowsAffected > 0
}

// GenerateToken generates a new JWT token for the given user ID.
//
// This function generates a new JWT token with the given user ID as the "userID" claim and a default expiration time of 24 hours.
// The token is signed using the secret key retrieved from the "APP_KEY" environment variable. If the environment variable is not set,
// a default secret key of "my-secret-key" is used. The function returns the generated token as a string and an error, if any.
//
// Example usage:
//
//	userID := 1234
//	token, err := GenerateToken(userID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Token:", token)
//
// Parameters:
//
//	userID (int): The ID of the user for whom the token is being generated.
//
// Returns:
//
//	string: The generated JWT token.
//	error: An error, if any, that occurred during the token generation process.
func GenerateToken(userID int) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// Get the secret key from an environment variable
	secretKey := []byte(Env("APP_KEY", "my-secret-key"))

	// Generate the signed token string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil // Return a jwt token
}

// DecodeRequestBody decodes the form data in the http request body and maps it to a struct
//
// This function takes an http request object and a struct to be decoded into, and returns an error
// if the decoding process fails. It first parses the request form data and then uses the
// github.com/gorilla/schema package to decode the form values into the given struct.
//
// Example usage:
//
//	var req MyRequestStruct
//	err := DecodeRequestBody(r, &req)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Parameters:
//   - r (*http.Request): the http request object containing the form data to be decoded
//   - v (interface{}): a pointer to the struct to be decoded into
//
// Returns:
//   - error: an error, if any, that occurred during the decoding process
func DecodeRequestBody(r *http.Request, v interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(v, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}

// DecodeJSONBody decodes the JSON request body into the provided destination object.
// It performs validation and error handling for the decoding process.
//
// Parameters:
//   - w: http.ResponseWriter - the response writer to send the error response if necessary
//   - r: *http.Request - the incoming HTTP request
//   - dst: interface{} - a pointer to the destination object where the JSON data will be decoded into
//
// Returns:
//   - error: an error indicating any issues that occurred during the decoding process
//   - nil if the decoding is successful
//   - MalformedRequest error if the request body or its contents are malformed
//   - Other decoding-related errors for specific error scenarios
//
// The function expects the Content-Type header of the request to be "application/json".
// If the header is missing or not set to "application/json", it returns a MalformedRequest error.
//
// The function limits the size of the request body to 1MB by using http.MaxBytesReader.
//
// The function uses the json.Decoder to decode the JSON data and disallows unknown fields.
// It handles various error scenarios during decoding, including syntax errors, unexpected EOF,
// invalid field values, unknown fields, empty body, and exceeding the size limit.
//
// Example usage:
//
//	type CreatePostRequest struct {
//	    Title string `json:"title" validate:"required,min=20"`
//	    Body  string `json:"body" validate:"required,min=40"`
//	}
//
//	func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
//	    var createReq CreatePostRequest
//	    if err := DecodeJSONBody(w, r, &createReq); err != nil {
//	        // Handle the error (e.g., return an error response)
//	        return
//	    }
//
//	    // Use the decoded data in further processing
//	}
//
// The function returns nil if the decoding is successful.
// Otherwise, it returns a MalformedRequest error or other decoding-related errors.
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &MalformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &MalformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &MalformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &MalformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &MalformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

// GenerateRandomString generates a random string of the given length using cryptographically secure random bytes
//
// This function takes an integer length and returns a random string of that length. The string is generated by
// generating random bytes using the rand package, and then encoding those bytes using base64 encoding. The resulting
// string is URL-safe, meaning it contains only URL-safe characters (i.e., it does not contain any '+' or '/' characters).
//
// Example usage:
//
//	randomString, err := GenerateRandomString(16)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Parameters:
//   - length (int): the length of the generated random string
//
// Returns:
//   - string: a random string of the given length
//   - error: an error, if any, that occurred during the generation process
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
