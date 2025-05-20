package customer_ticket

// CustomerCommentCommand represents a comment on a customer ticket
type CustomerCommentCommand struct {
	Content      *string `json:"content" nameFa:"محتوا" validate:"required_text=1 2000"`
	RespondentID *int64  `json:"respondentId" nameFa:"شناسه پاسخگو" validate:"required,gt=0"`
}
