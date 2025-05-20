package ticket

// CreateTicketCommand represents a command to create a new ticket
type CreateTicketCommand struct {
	Title    *string             `json:"title" validate:"required_text=1,200"`
	Category *TicketCategoryEnum `json:"product_category" validate:"required,enum"`
	Priority *TicketPriorityEnum `json:"priority" validate:"required,enum"`
	Comment  *CommentCommand     `json:"comment" validate:"required"`
	MediaIDs []int64             `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
}

// ReplayTicketCommand represents a command to reply to a ticket
type ReplayTicketCommand struct {
	ID         *int64              `json:"id" validate:"required,gt=0"`
	Status     *TicketStatusEnum   `json:"status" validate:"required,enum"`
	Category   *TicketCategoryEnum `json:"product_category" validate:"required,enum"`
	AssignedTo *int64              `json:"assignedTo" validate:"required,gt=0"`
	Priority   *TicketPriorityEnum `json:"priority" validate:"required,enum"`
	Comment    *CommentCommand     `json:"comment" validate:"required"`
	MediaIDs   []int64             `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
}

// AdminReplayTicketCommand represents a command for an admin to reply to a ticket
type AdminReplayTicketCommand struct {
	ID         *int64              `json:"id" validate:"required,gt=0"`
	Status     *TicketStatusEnum   `json:"status" validate:"required,enum"`
	Category   *TicketCategoryEnum `json:"product_category" validate:"required,enum"`
	AssignedTo *int64              `json:"assignedTo,omitempty" validate:"omitempty,gt=0"`
	Priority   *TicketPriorityEnum `json:"priority" validate:"required,enum"`
	Comment    *CommentCommand     `json:"comment" validate:"required"`
	MediaIDs   []int64             `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false"`
}
