package user

// ResponseStatus represents a response status for payment verification
type ResponseStatus struct {
	IsSuccess bool   `json:"isSuccess" nameFa:"موفقیت"`
	Message   string `json:"message" nameFa:"پیام"`
}

// VerifyPaymentResponseDto represents a response for payment verification
type VerifyPaymentResponseDto struct {
	ResponseStatus ResponseStatus `json:"responseStatus"`
}
