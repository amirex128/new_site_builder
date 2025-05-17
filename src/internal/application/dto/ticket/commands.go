package ticket

// CreateTicketCommand represents a command to create a new ticket
type CreateTicketCommand struct {
	Title    *string             `json:"title" validate:"required,max=200" error:"required=عنوان الزامی است|max=عنوان نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Category *TicketCategoryEnum `json:"product_category" validate:"required" error:"required=دسته‌بندی الزامی است"`
	Priority *TicketPriorityEnum `json:"priority" validate:"required" error:"required=اولویت الزامی است"`
	Comment  *CommentCommand     `json:"comment" validate:"required" error:"required=نظر الزامی است"`
	MediaIDs []int64             `json:"mediaIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های رسانه باید بزرگتر از 0 باشند"`
}

// ReplayTicketCommand represents a command to reply to a ticket
type ReplayTicketCommand struct {
	ID         *int64              `json:"id" validate:"required,gt=0" error:"required=شناسه تیکت الزامی است|gt=شناسه تیکت باید بزرگتر از 0 باشد"`
	Status     *TicketStatusEnum   `json:"status" validate:"required" error:"required=وضعیت الزامی است"`
	Category   *TicketCategoryEnum `json:"product_category" validate:"required" error:"required=دسته‌بندی الزامی است"`
	AssignedTo *int64              `json:"assignedTo" validate:"required,gt=0" error:"required=شناسه کاربر تخصیص داده شده الزامی است|gt=شناسه کاربر تخصیص داده شده باید بزرگتر از 0 باشد"`
	Priority   *TicketPriorityEnum `json:"priority" validate:"required" error:"required=اولویت الزامی است"`
	Comment    *CommentCommand     `json:"comment" validate:"required" error:"required=نظر الزامی است"`
	MediaIDs   []int64             `json:"mediaIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های رسانه باید بزرگتر از 0 باشند"`
}

// AdminReplayTicketCommand represents a command for an admin to reply to a ticket
type AdminReplayTicketCommand struct {
	ID         *int64              `json:"id" validate:"required,gt=0" error:"required=شناسه تیکت الزامی است|gt=شناسه تیکت باید بزرگتر از 0 باشد"`
	Status     *TicketStatusEnum   `json:"status" validate:"required" error:"required=وضعیت الزامی است"`
	Category   *TicketCategoryEnum `json:"product_category" validate:"required" error:"required=دسته‌بندی الزامی است"`
	AssignedTo *int64              `json:"assignedTo,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه کاربر تخصیص داده شده باید بزرگتر از 0 باشد"` // Note: Changed to optional based on the original validation
	Priority   *TicketPriorityEnum `json:"priority" validate:"required" error:"required=اولویت الزامی است"`
	Comment    *CommentCommand     `json:"comment" validate:"required" error:"required=نظر الزامی است"`
	MediaIDs   []int64             `json:"mediaIds,omitempty" validate:"omitempty,dive,gt=0" error:"gt=شناسه‌های رسانه باید بزرگتر از 0 باشند"`
}
