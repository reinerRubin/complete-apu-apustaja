package handler

type ErrorResponse struct {
	Msg string
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Msg: msg}
}
