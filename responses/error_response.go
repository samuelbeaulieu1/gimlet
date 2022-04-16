package responses

type ErrorResponse struct {
	*RequestResponse
	Data *RequestResponseMessage `json:"data"`
}

type Error interface {
	Error() string
}

func NewError(message string) *ErrorResponse {
	return &ErrorResponse{
		&RequestResponse{
			Status: "error",
		},
		&RequestResponseMessage{
			message,
		},
	}
}

func (err *ErrorResponse) Error() string {
	return err.Data.Message
}
