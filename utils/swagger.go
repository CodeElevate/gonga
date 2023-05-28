package utils

// type APIResponse struct {
// 	Type    string                 `json:"type,omitempty"`
// 	Message string                 `json:"message,omitempty"`
// 	Meta    map[string]interface{} `json:"meta,omitempty"`
// 	Data    interface{}            `json:"data,omitempty"`
// 	Errors  []APIError             `json:"errors,omitempty"`
// }

// type APIError struct {
// 	Code       int                    `json:"code,omitempty" example:"400"`
// 	Message    string                 `json:"message,omitempty" example:"Bad Request"`
// 	Field      string                 `json:"field,omitempty" example:"username"`
// 	Details    map[string]interface{} `json:"details,omitempty"`
// 	Critical   bool                   `json:"critical,omitempty" example:"true"`
// 	Suggestion string                 `json:"suggestion,omitempty" example:"Try again later."`
// }
type SwaggerPagination struct {
	Type    string                 `json:"type,omitempty"`
	Message string                 `json:"message,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Data    interface{}            `json:"data,omitempty"`
}

type SwaggerSuccessResponse struct {
	Type    string      `json:"type,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type SwaggerErrorResponse struct {
	Type    string                 `json:"type,omitempty"`
	Message string                 `json:"message,omitempty"`
	Errors  []APIError             `json:"errors,omitempty"`
}