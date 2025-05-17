package customer_ticket

// CustomerCommentCommand represents a comment on a customer ticket
type CustomerCommentCommand struct {
	Content      *string `json:"content" validate:"required,max=2000" error:"required=محتوای نظر الزامی است|max=محتوای نظر نمی‌تواند بیشتر از 2000 کاراکتر باشد"`
	RespondentID *int64  `json:"respondentId" validate:"required,gt=0" error:"required=شناسه پاسخ دهنده الزامی است|gt=شناسه پاسخ دهنده باید بزرگتر از 0 باشد"`
}

// CustomerTicket represents a customer support ticket
type CustomerTicket struct {
	ID          *int64                      `json:"id,omitempty" validate:"omitempty" error:""`
	Title       *string                     `json:"title" validate:"required,max=200" error:"required=عنوان الزامی است|max=عنوان نمی‌تواند بیشتر از 200 کاراکتر باشد"`
	OwnerUserID *int64                      `json:"ownerUserId" validate:"required,gt=0" error:"required=شناسه کاربر صاحب تیکت الزامی است|gt=شناسه کاربر صاحب تیکت باید بزرگتر از 0 باشد"`
	Status      *CustomerTicketStatusEnum   `json:"status" validate:"required" error:"required=وضعیت الزامی است"`
	Category    *CustomerTicketCategoryEnum `json:"product_category" validate:"required" error:"required=دسته‌بندی الزامی است"`
	AssignedTo  *int64                      `json:"assignedTo,omitempty" validate:"omitempty,gt=0" error:"gt=شناسه کاربر تخصیص داده شده باید بزرگتر از 0 باشد"`
	Priority    *CustomerTicketPriorityEnum `json:"priority" validate:"required" error:"required=اولویت الزامی است"`
}
