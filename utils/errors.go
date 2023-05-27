package utils

import "net/http"

func HandleError(w http.ResponseWriter, err error, statusCode int, message ...string) {
	response := APIResponse{
		Type: "error",
		Errors: []APIError{
			{
				Message: err.Error(),
				Code:    statusCode,
			},
		},
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	JSONResponse(w, statusCode, response)
}
