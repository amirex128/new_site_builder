package plan

import "github.com/amirex128/new_site_builder/src/internal/domain/enums"

// CreatePlanCommand represents a command to create a new plan
type CreatePlanCommand struct {
	Name             *string                 `json:"name" validate:"required_text=1 100"`
	ShowStatus       *string                 `json:"showStatus" validate:"required_text=1 50"`
	Description      *string                 `json:"description" validate:"required_text=1 500"`
	Price            *int                    `json:"price" validate:"required,min=0"`
	DiscountType     *enums.DiscountTypeEnum `json:"discountType,omitempty" validate:"required_with=Discount,enum_optional"`
	Discount         *int                    `json:"discount,omitempty" validate:"omitempty,min=0,max=100"`
	Duration         *int                    `json:"duration" validate:"required,min=1"`
	Feature          *string                 `json:"feature,omitempty" validate:"optional_text=1 1000"`
	SmsCredits       *int                    `json:"smsCredits" validate:"required,min=0"`
	EmailCredits     *int                    `json:"emailCredits" validate:"required,min=0"`
	StorageMbCredits *int                    `json:"storageMbCredits" validate:"required,min=0"`
	AiCredits        *int                    `json:"aiCredits" validate:"required,min=0"`
	AiImageCredits   *int                    `json:"aiImageCredits" validate:"required,min=0"`
}

// DeletePlanCommand represents a command to delete a plan
type DeletePlanCommand struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}

// UpdatePlanCommand represents a command to update a plan
type UpdatePlanCommand struct {
	ID               *int64                  `json:"id" validate:"required,gt=0"`
	Name             *string                 `json:"name,omitempty" validate:"optional_text=1 100"`
	ShowStatus       *string                 `json:"showStatus,omitempty" validate:"optional_text=1 50"`
	Description      *string                 `json:"description,omitempty" validate:"optional_text=1 500"`
	Price            *int                    `json:"price,omitempty" validate:"omitempty,min=0"`
	DiscountType     *enums.DiscountTypeEnum `json:"discountType,omitempty" validate:"required_with=Discount,enum_optional"`
	Discount         *int                    `json:"discount,omitempty" validate:"omitempty,min=0,max=100"`
	Duration         *int                    `json:"duration,omitempty" validate:"omitempty,min=1"`
	Feature          *string                 `json:"feature,omitempty" validate:"optional_text=1 1000"`
	SmsCredits       *int                    `json:"smsCredits,omitempty" validate:"omitempty,min=0"`
	EmailCredits     *int                    `json:"emailCredits,omitempty" validate:"omitempty,min=0"`
	StorageMbCredits *int                    `json:"storageMbCredits,omitempty" validate:"omitempty,min=0"`
	AiCredits        *int                    `json:"aiCredits,omitempty" validate:"omitempty,min=0"`
	AiImageCredits   *int                    `json:"aiImageCredits,omitempty" validate:"omitempty,min=0"`
}
