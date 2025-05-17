package ticket

// CommentCommand represents a comment on a ticket
type CommentCommand struct {
	Content      *string `json:"content" validate:"required_text=1,2000"`
	RespondentID *int64  `json:"respondentId" validate:"required,gt=0"`
}

// Ticket represents a support ticket
type Ticket struct {
	ID         *int64              `json:"id,omitempty" validate:"omitempty"`
	Title      *string             `json:"title" validate:"required_text=1,200"`
	Status     *TicketStatusEnum   `json:"status" validate:"required,enum"`
	Category   *TicketCategoryEnum `json:"product_category" validate:"required,enum"`
	AssignedTo *int64              `json:"assignedTo,omitempty" validate:"omitempty,gt=0"`
	Priority   *TicketPriorityEnum `json:"priority" validate:"required,enum"`
	Media      []interface{}       `json:"media,omitempty" validate:"omitempty"`
}
