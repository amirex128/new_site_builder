package customer_ticket

import (
	enums2 "github.com/amirex128/new_site_builder/internal/domain/enums"
)

// CreateCustomerTicketCommand represents a command to create a new customer ticket
type CreateCustomerTicketCommand struct {
	Title       *string                    `json:"title" nameFa:"عنوان" validate:"required_text=1 200"`
	OwnerUserID *int64                     `json:"ownerUserId" nameFa:"شناسه کاربر مالک" validate:"required,gt=0"`
	Category    *enums2.TicketCategoryEnum `json:"product_category" nameFa:"دسته بندی" validate:"required,enum"`
	Priority    *enums2.TicketPriorityEnum `json:"priority" nameFa:"اولویت" validate:"required,enum"`
	Comment     *CustomerCommentCommand    `json:"comment" nameFa:"نظر" validate:"required"`
	MediaIDs    []int64                    `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 100 1 0 false"`
}

// ReplayCustomerTicketCommand represents a command to reply to a customer ticket
type ReplayCustomerTicketCommand struct {
	ID         *int64                     `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
	Status     *enums2.TicketStatusEnum   `json:"status" nameFa:"وضعیت" validate:"required,enum"`
	Category   *enums2.TicketCategoryEnum `json:"product_category" nameFa:"دسته بندی" validate:"required,enum"`
	AssignedTo *int64                     `json:"assignedTo" nameFa:"شناسه کاربر تخصیص داده شده" validate:"required,gt=0"`
	Priority   *enums2.TicketPriorityEnum `json:"priority" nameFa:"اولویت" validate:"required,enum"`
	Comment    *CustomerCommentCommand    `json:"comment" nameFa:"نظر" validate:"required"`
	MediaIDs   []int64                    `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 100 1 0 false"`
}

// AdminReplayCustomerTicketCommand represents a command for an admin to reply to a customer ticket
type AdminReplayCustomerTicketCommand struct {
	ID         *int64                     `json:"id" nameFa:"شناسه" validate:"required,gt=0"`
	Status     *enums2.TicketStatusEnum   `json:"status" nameFa:"وضعیت" validate:"required,enum"`
	Category   *enums2.TicketCategoryEnum `json:"product_category" nameFa:"دسته بندی" validate:"required,enum"`
	AssignedTo *int64                     `json:"assignedTo" nameFa:"شناسه کاربر تخصیص داده شده" validate:"required,gt=0"`
	Priority   *enums2.TicketPriorityEnum `json:"priority" nameFa:"اولویت" validate:"required,enum"`
	Comment    *CustomerCommentCommand    `json:"comment" nameFa:"نظر" validate:"required"`
	MediaIDs   []int64                    `json:"mediaIds,omitempty" nameFa:"شناسه های مدیا" validate:"array_number_optional=0 100 1 0 false"`
}
