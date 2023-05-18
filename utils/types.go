package utils

type PaginatedResult struct {
    Page         int         `json:"page"`
    PerPage      int         `json:"perPage"`
    TotalRecords int         `json:"totalRecords"`
    TotalPages   int         `json:"totalPages"`
    Items        interface{} `json:"items"`
    Remaining    int         `json:"remaining"`
}

type MalformedRequest struct {
	status int
	msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.msg
}

func (mr *MalformedRequest) Status() int {
	return mr.status
}

type ContextKey string
