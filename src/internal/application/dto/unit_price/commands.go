package unit_price

import (
	"github.com/amirex128/new_site_builder/src/internal/application/dto/user"
)

// UpdateUnitPriceCommand represents a command to update a unit price
type UpdateUnitPriceCommand struct {
	ID           *int64            `json:"id" validate:"required,gt=0" error:"required=واحد قیمت الزامی است|gt=شناسه واحد قیمت باید بزرگتر از 0 باشد"`
	Name         *string           `json:"name,omitempty" validate:"omitempty,max=100" error:"max=نام نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Price        *int              `json:"price,omitempty" error:""`
	HasDay       *bool             `json:"hasDay" validate:"required" error:"required=وضعیت روز الزامی است"`
	DiscountType *DiscountTypeEnum `json:"discountType,omitempty" validate:"required_with=Discount" error:"required_with=نوع تخفیف الزامی است"`
	Discount     *int              `json:"discount,omitempty" validate:"omitempty,min=0,max=100" error:"min=تخفیف نمی‌تواند کمتر از 0 باشد|max=تخفیف نمی‌تواند بیشتر از 100 باشد"`
}
