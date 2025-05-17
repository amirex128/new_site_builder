package plan

// CreatePlanCommand represents a command to create a new plan
type CreatePlanCommand struct {
	Name           *string           `json:"name" validate:"required_text=1,100"`
	Description    *string           `json:"description,omitempty" validate:"optional_text=1,1000"`
	Price          *int64            `json:"price" validate:"required,min=0"`
	DiscountType   *DiscountTypeEnum `json:"discountType,omitempty" validate:"enum_optional"`
	Discount       *int64            `json:"discount,omitempty" validate:"omitempty,min=0"`
	Duration       *int              `json:"duration" validate:"required,min=1,max=1000"`
	Feature        *string           `json:"feature,omitempty" validate:"optional_text=1,500"`
	SmsCredits     *int              `json:"smsCredits,omitempty" validate:"omitempty,min=0,max=1000"`
	EmailCredits   *int              `json:"emailCredits,omitempty" validate:"omitempty,min=0,max=1000"`
	StorageCredits *int              `json:"storageCredits,omitempty" validate:"omitempty,min=0,max=1000"`
	AiCredits      *int              `json:"aiCredits,omitempty" validate:"omitempty,min=0,max=1000"`
	AiImageCredits *int              `json:"aiImageCredits,omitempty" validate:"omitempty,min=0,max=1000"`
}

// DeletePlanCommand represents a command to delete a plan
type DeletePlanCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}

// UpdatePlanCommand represents a command to update a plan
type UpdatePlanCommand struct {
	ID             *int64            `json:"id" validate:"required,gt=0"`
	Name           *string           `json:"name,omitempty" validate:"optional_text=1,100"`
	Description    *string           `json:"description,omitempty" validate:"optional_text=1,1000"`
	Price          *int64            `json:"price,omitempty" validate:"omitempty,min=0"`
	DiscountType   *DiscountTypeEnum `json:"discountType,omitempty" validate:"enum_optional"`
	Discount       *int64            `json:"discount,omitempty" validate:"omitempty,min=0"`
	Duration       *int              `json:"duration,omitempty" validate:"omitempty,min=1"`
	Feature        *string           `json:"feature,omitempty" validate:"optional_text=1,500"`
	SmsCredits     *int              `json:"smsCredits,omitempty" validate:"omitempty,min=0"`
	EmailCredits   *int              `json:"emailCredits,omitempty" validate:"omitempty,min=0"`
	StorageCredits *int              `json:"storageCredits,omitempty" validate:"omitempty,min=0"`
	AiCredits      *int              `json:"aiCredits,omitempty" validate:"omitempty,min=0"`
	AiImageCredits *int              `json:"aiImageCredits,omitempty" validate:"omitempty,min=0"`
}
