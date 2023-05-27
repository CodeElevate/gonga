package utils




// type APIResponse struct {
// 	Type    string                 `json:"type,omitempty"`
// 	Message string                 `json:"message,omitempty"`
// 	Meta    map[string]interface{} `json:"meta,omitempty"`
// 	Data    interface{}            `json:"data,omitempty"`
// 	Errors  []APIError             `json:"errors,omitempty"`
// }

// type APIError struct {
// 	Code       int                    `json:"code,omitempty"`
// 	Message    string                 `json:"message,omitempty"`
// 	Field      string                 `json:"field,omitempty"`
// 	Details    map[string]interface{} `json:"details,omitempty"`
// 	Critical   bool                   `json:"critical,omitempty"`
// 	Suggestion string                 `json:"suggestion,omitempty"`
// }
type SwaggerPagination struct{
	Type    string                 `json:"type,omitempty"`
	Message string                 `json:"message,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Data    interface{}            `json:"data,omitempty"`
}

type SwaggerSuccessResponse struct {
	Type string                 `json:"type,omitempty"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}
