package ticket

// CommentCommand represents a comment on a ticket
type CommentCommand struct {
	Content      *string `json:"content" validate:"required,max=2000" error:"required=محتوای نظر الزامی است|max=محتوای نظر نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	RespondentID *int64  `json:"respondentId" validate:"required,gt=0" error:"required=شناسه پاسخ دهنده الزامی است|gt=شناسه پاسخ دهنده باید بزرگتر از 0 باشد"`
}

// Ticket represents a support ticket
type Ticket struct {
	ID         *int64              `json:"id,omitempty" validate:"omitempty" error:""`
	Title      *string             `json:"title" validate:"required,max=200" error:"required=عنوان الزامی است|max=عنوان نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	Status     *TicketStatusEnum   `json:"status" validate:"required" error:"required=وضعیت الزامی است"`
	Category   *TicketCategoryEnum `json:"product_category" validate:"required" error:"required=دسته‌بندی الزامی است"`
	AssignedTo *int64              `json:"assignedTo,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه کاربر تخصیص داده شده باید بزرگتر از 0 باشد"`
	Priority   *TicketPriorityEnum `json:"priority" validate:"required" error:"required=اولویت الزامی است"`
	Media      []interface{}       `json:"media,omitempty" validate:"omitempty" error:""`
}
