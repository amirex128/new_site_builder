package unit_price

import "github.com/amirex128/new_site_builder/src/internal/domain/enums"

// UpdateUnitPriceCommand represents a command to update a unit price
type UpdateUnitPriceCommand struct {
	ID           *int64                  `json:"id" validate:"required,gt=0" nameFa:"شناسه واحد قیمت"`
	Name         *string                 `json:"name,omitempty" validate:"optional_text=1 100" nameFa:"نام واحد قیمت"`
	Price        *int                    `json:"price,omitempty" nameFa:"قیمت"`
	HasDay       *bool                   `json:"hasDay" validate:"required_bool" nameFa:"دارای روز"`
	DiscountType *enums.DiscountTypeEnum `json:"discountType,omitempty" validate:"required_with=Discount,enum_optional" nameFa:"نوع تخفیف"`
	Discount     *int                    `json:"discount,omitempty" validate:"omitempty,min=0,max=100" nameFa:"میزان تخفیف"`
}
