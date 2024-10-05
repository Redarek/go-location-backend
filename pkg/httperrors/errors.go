package httperrors

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func NewErrorResponse(code int, errorType string, message string, details interface{}) *ErrorResponse {
	return &ErrorResponse{
		Error:   errorType,
		Message: message,
		Code:    code,
		Details: details,
	}
}
