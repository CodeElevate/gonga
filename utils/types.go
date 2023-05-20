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

type ContextKey string
