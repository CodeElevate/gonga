package utils

type MalformedRequest struct {
	status int
	msg    string
}
type ControllerResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (mr *MalformedRequest) Error() string {
	return mr.msg
}

func (mr *MalformedRequest) Status() int {
	return mr.status
}

type APIResponse struct {
	Type    string                 `json:"type,omitempty"`
	Message string                 `json:"message,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Data    interface{}            `json:"data,omitempty"`
	Errors  []APIError             `json:"errors,omitempty"`
}

type APIError struct {
	Code       int                    `json:"code,omitempty" example:"400"`
	Message    string                 `json:"message,omitempty" example:"Bad Request"`
	Field      string                 `json:"field,omitempty" example:"username"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Critical   bool                   `json:"critical,omitempty" example:"true"`
	Suggestion string                 `json:"suggestion,omitempty" example:"Try again later."`
}

type ContextKey string
