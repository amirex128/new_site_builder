package customer_ticket

// CustomerCommentCommand represents a comment on a customer ticket
type CustomerCommentCommand struct {
	Content      *string `json:"content" validate:"required_text=1,2000"`
	RespondentID *int64  `json:"respondentId" validate:"required,gt=0"`
}

// CustomerTicket represents a customer support ticket
type CustomerTicket struct {
	ID          *int64                      `json:"id,omitempty" validate:"omitempty"`
	Title       *string                     `json:"title" validate:"required_text=1,200"`
	OwnerUserID *int64                      `json:"ownerUserId" validate:"required,gt=0"`
	Status      *CustomerTicketStatusEnum   `json:"status" validate:"required,enum"`
	Category    *CustomerTicketCategoryEnum `json:"product_category" validate:"required,enum"`
	AssignedTo  *int64                      `json:"assignedTo,omitempty" validate:"omitempty,gt=0"`
	Priority    *CustomerTicketPriorityEnum `json:"priority" validate:"required,enum"`
}
