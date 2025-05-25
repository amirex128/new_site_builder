package resp

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
	Data    any       `json:"data"`
}
type Response struct {
	Message string
	Type    ResponseType
	Data    any
}

func NewError(typ ErrorType, msg string) *Error {
	return &Error{
		Message: msg,
		Type:    typ,
	}
}

func NewErrorData(typ ErrorType, data any, msg string) *Error {
	return &Error{
		Message: msg,
		Type:    typ,
		Data:    data,
	}
}
func NewResponseData(typ ResponseType, data any, msg string) *Response {
	return &Response{
		Message: msg,
		Type:    typ,
		Data:    data,
	}
}
func NewResponse(typ ResponseType, msg string) *Response {
	return &Response{
		Message: msg,
		Type:    typ}
}

func (e *Error) Error() string {
	return e.Message
}
