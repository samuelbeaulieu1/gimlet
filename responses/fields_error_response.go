package responses

import (
	"strings"
)

type FieldErrorPayload struct {
	Messages []string `json:"messages"`
	Fields   []string `json:"fields"`
}

type FieldErrorResponse struct {
	*RequestResponse
	Data *FieldErrorPayload `json:"data"`
}

func NewFieldsError(messages []string, fields []string) *FieldErrorResponse {
	fieldErr := &FieldErrorResponse{
		&RequestResponse{
			Status: "fail",
		},
		&FieldErrorPayload{
			messages,
			fields,
		},
	}

	return fieldErr
}

func (err *FieldErrorResponse) Error() string {
	return strings.Join(err.Data.Messages, "\n")
}
