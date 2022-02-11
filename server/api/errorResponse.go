package api

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(err error) *errorResponse {
	return &errorResponse{err.Error()}
}
