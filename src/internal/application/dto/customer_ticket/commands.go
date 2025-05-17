package customer_ticket

// CreateCustomerTicketCommand represents a command to create a new customer ticket
type CreateCustomerTicketCommand struct {
	Title       *string                     `json:"title" validate:"required_text=1,200"`
	OwnerUserID *int64                      `json:"ownerUserId" validate:"required,gt=0"`
	Category    *CustomerTicketCategoryEnum `json:"product_category" validate:"required,enum"`
	Priority    *CustomerTicketPriorityEnum `json:"priority" validate:"required,enum"`
	Comment     *CustomerCommentCommand     `json:"comment" validate:"required"`
	MediaIDs    []int64                     `json:"mediaIds,omitempty" validate:"array_number_optional=0,100,1,0,false"`
}

// ReplayCustomerTicketCommand represents a command to reply to a customer ticket
type ReplayCustomerTicketCommand struct {
	ID         *int64                      `json:"id" validate:"required,gt=0"`
	Status     *CustomerTicketStatusEnum   `json:"status" validate:"required,enum"`
	Category   *CustomerTicketCategoryEnum `json:"product_category" validate:"required,enum"`
	AssignedTo *int64                      `json:"assignedTo" validate:"required,gt=0"`
	Priority   *CustomerTicketPriorityEnum `json:"priority" validate:"required,enum"`
	Comment    *CustomerCommentCommand     `json:"comment" validate:"required"`
	MediaIDs   []int64                     `json:"mediaIds,omitempty" validate:"array_number_optional=0,100,1,0,false"`
}

// AdminReplayCustomerTicketCommand represents a command for an admin to reply to a customer ticket
type AdminReplayCustomerTicketCommand struct {
	ID         *int64                      `json:"id" validate:"required,gt=0"`
	Status     *CustomerTicketStatusEnum   `json:"status" validate:"required,enum"`
	Category   *CustomerTicketCategoryEnum `json:"product_category" validate:"required,enum"`
	AssignedTo *int64                      `json:"assignedTo" validate:"required,gt=0"`
	Priority   *CustomerTicketPriorityEnum `json:"priority" validate:"required,enum"`
	Comment    *CustomerCommentCommand     `json:"comment" validate:"required"`
	MediaIDs   []int64                     `json:"mediaIds,omitempty" validate:"array_number_optional=0,100,1,0,false"`
}
