package resp

import "fmt"

type ErrorType string
type ResponseType string
type Data map[string]any

const (
	NotFound     ErrorType = "not_found"
	Internal     ErrorType = "internal"
	Unauthorized ErrorType = "unauthorized"
	Validation   ErrorType = "validation"
	BadRequest   ErrorType = "bad_request"
)
const (
	Success   ResponseType = "success"
	Created   ResponseType = "created"
	Updated   ResponseType = "updated"
	Deleted   ResponseType = "deleted"
	Retrieved ResponseType = "retrieved"
	NoContent ResponseType = "no_content"
)

type Error struct {
	Message string    `json:"message"`
	Type    ErrorType `json:"type"`
	Data    Data      `json:"data"`
}
type Response struct {
	Message string
	Type    ResponseType
	Data    Data
}

func NewError(typ ErrorType, msg string, format ...any) *Error {
	return &Error{
		Message: fmt.Sprintf(msg, format...),
		Type:    typ,
	}
}

func NewErrorData(typ ErrorType, data Data, msg string, format ...any) *Error {
	return &Error{
		Message: fmt.Sprintf(msg, format...),
		Type:    typ,
		Data:    data,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func NewResponseData(typ ResponseType, data Data, msg string, format ...any) *Response {
	return &Response{
		Message: fmt.Sprintf(msg, format...),
		Type:    typ,
		Data:    data,
	}
}
func NewResponse(typ ResponseType, msg string, format ...any) *Response {
	return &Response{
		Message: fmt.Sprintf(msg, format...),
		Type:    typ}
}
