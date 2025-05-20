package unit_price

import "github.com/amirex128/new_site_builder/src/internal/domain/enums"

// UpdateUnitPriceCommand represents a command to update a unit price
type UpdateUnitPriceCommand struct {
	ID           *int64                  `json:"id" validate:"required,gt=0"`
	Name         *string                 `json:"name,omitempty" validate:"optional_text=1 100"`
	Price        *int                    `json:"price,omitempty"`
	HasDay       *bool                   `json:"hasDay" validate:"required_bool"`
	DiscountType *enums.DiscountTypeEnum `json:"discountType,omitempty" validate:"required_with=Discount,enum_optional"`
	Discount     *int                    `json:"discount,omitempty" validate:"omitempty,min=0,max=100"`
}
