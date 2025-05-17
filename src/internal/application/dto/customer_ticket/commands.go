package customer_ticket

// CreateCustomerTicketCommand represents a command to create a new customer ticket
type CreateCustomerTicketCommand struct {
	Title       *string                     `json:"title" validate:"required,max=200" error:"required=عنوان الزامی است|max=عنوان نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	OwnerUserID *int64                      `json:"ownerUserId" validate:"required,gt=0" error:"required=شناسه کاربر صاحب تیکت الزامی است|gt=شناسه کاربر صاحب تیکت باید بزرگتر از 0 باشد"`
	Category    *CustomerTicketCategoryEnum `json:"product_category" validate:"required" error:"required=دسته‌بندی الزامی است"`
	Priority    *CustomerTicketPriorityEnum `json:"priority" validate:"required" error:"required=اولویت الزامی است"`
	Comment     *CustomerCommentCommand     `json:"comment" validate:"required" error:"required=نظر الزامی است"`
	MediaIDs    []int64                     `json:"mediaIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های رسانه باید بزرگتر از 0 باشند"`
}

// ReplayCustomerTicketCommand represents a command to reply to a customer ticket
type ReplayCustomerTicketCommand struct {
	ID         *int64                      `json:"id" validate:"required,gt=0" error:"required=شناسه تیکت الزامی است|gt=شناسه تیکت باید بزرگتر از 0 باشد"`
	Status     *CustomerTicketStatusEnum   `json:"status" validate:"required" error:"required=وضعیت الزامی است"`
	Category   *CustomerTicketCategoryEnum `json:"product_category" validate:"required" error:"required=دسته‌بندی الزامی است"`
	AssignedTo *int64                      `json:"assignedTo" validate:"required,gt=0" error:"required=شناسه کاربر تخصیص داده شده الزامی است|gt=شناسه کاربر تخصیص داده شده باید بزرگتر از 0 باشد"`
	Priority   *CustomerTicketPriorityEnum `json:"priority" validate:"required" error:"required=اولویت الزامی است"`
	Comment    *CustomerCommentCommand     `json:"comment" validate:"required" error:"required=نظر الزامی است"`
	MediaIDs   []int64                     `json:"mediaIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های رسانه باید بزرگتر از 0 باشند"`
}

// AdminReplayCustomerTicketCommand represents a command for an admin to reply to a customer ticket
type AdminReplayCustomerTicketCommand struct {
	ID         *int64                      `json:"id" validate:"required,gt=0" error:"required=شناسه تیکت الزامی است|gt=شناسه تیکت باید بزرگتر از 0 باشد"`
	Status     *CustomerTicketStatusEnum   `json:"status" validate:"required" error:"required=وضعیت الزامی است"`
	Category   *CustomerTicketCategoryEnum `json:"product_category" validate:"required" error:"required=دسته‌بندی الزامی است"`
	AssignedTo *int64                      `json:"assignedTo" validate:"required,gt=0" error:"required=شناسه کاربر تخصیص داده شده الزامی است|gt=شناسه کاربر تخصیص داده شده باید بزرگتر از 0 باشد"`
	Priority   *CustomerTicketPriorityEnum `json:"priority" validate:"required" error:"required=اولویت الزامی است"`
	Comment    *CustomerCommentCommand     `json:"comment" validate:"required" error:"required=نظر الزامی است"`
	MediaIDs   []int64                     `json:"mediaIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های رسانه باید بزرگتر از 0 باشند"`
}
