package ticket

import (
	enums2 "github.com/amirex128/new_site_builder/internal/domain/enums"
)

// CreateTicketCommand represents a command to create a new ticket
type CreateTicketCommand struct {
	Title    *string                    `json:"title" validate:"required_text=1,200" nameFa:"عنوان"`
	Category *enums2.TicketCategoryEnum `json:"product_category" validate:"required,enum" nameFa:"نوع محصول"`
	Priority *enums2.TicketPriorityEnum `json:"priority" validate:"required,enum" nameFa:"اولویت"`
	Comment  *CommentCommand            `json:"comment" validate:"required" nameFa:"نظر"`
	MediaIDs []int64                    `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false" nameFa:"شناسه مدیا"`
}

// ReplayTicketCommand represents a command to reply to a ticket
type ReplayTicketCommand struct {
	ID         *int64                     `json:"id" validate:"required,gt=0" nameFa:"شناسه"`
	Status     *enums2.TicketStatusEnum   `json:"status" validate:"required,enum" nameFa:"وضعیت"`
	Category   *enums2.TicketCategoryEnum `json:"product_category" validate:"required,enum" nameFa:"نوع محصول"`
	AssignedTo *int64                     `json:"assignedTo" validate:"required,gt=0" nameFa:"اختصاص داده شده به"`
	Priority   *enums2.TicketPriorityEnum `json:"priority" validate:"required,enum" nameFa:"اولویت"`
	Comment    *CommentCommand            `json:"comment" validate:"required" nameFa:"نظر"`
	MediaIDs   []int64                    `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false" nameFa:"شناسه مدیا"`
}

// AdminReplayTicketCommand represents a command for an admin to reply to a ticket
type AdminReplayTicketCommand struct {
	ID         *int64                     `json:"id" validate:"required,gt=0" nameFa:"شناسه"`
	Status     *enums2.TicketStatusEnum   `json:"status" validate:"required,enum" nameFa:"وضعیت"`
	Category   *enums2.TicketCategoryEnum `json:"product_category" validate:"required,enum" nameFa:"نوع محصول"`
	AssignedTo *int64                     `json:"assignedTo,omitempty" validate:"omitempty,gt=0" nameFa:"اختصاص داده شده به"`
	Priority   *enums2.TicketPriorityEnum `json:"priority" validate:"required,enum" nameFa:"اولویت"`
	Comment    *CommentCommand            `json:"comment" validate:"required" nameFa:"نظر"`
	MediaIDs   []int64                    `json:"mediaIds,omitempty" validate:"array_number_optional=0 100 1 0 false" nameFa:"شناسه مدیا"`
}
