package ticket

// CommentCommand represents a comment on a ticket
type CommentCommand struct {
	Content      *string       `json:"content" validate:"required_text=1 2000" nameFa:"محتوا"`
	RespondentID *int64        `json:"respondentId" validate:"required,gt=0" nameFa:"شناسه پاسخگو"`
	Media        []interface{} `json:"media,omitempty" validate:"omitempty" nameFa:"رسانه"`
}
