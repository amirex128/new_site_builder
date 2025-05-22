package nerror

import "fmt"

type ErrorType string

const (
	NotFound     ErrorType = "not_found"
	Internal     ErrorType = "internal"
	Unauthorized ErrorType = "unauthorized"
	BadRequest   ErrorType = "bad_request"
)

type NError struct {
	Message string    `json:"message"`
	Type    ErrorType `json:"type"`
}

func NewNError(typ ErrorType, msg string, format ...any) *NError {
	return &NError{
		Message: fmt.Sprintf(msg, format...),
		Type:    typ,
	}
}
func (e *NError) Error() string {
	return e.Message
}

func (e *NError) ErrorType() ErrorType {
	return e.Type
}
