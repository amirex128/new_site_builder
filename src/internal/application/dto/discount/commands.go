package discount

import (
	"time"

	"github.com/amirex128/new_site_builder/src/internal/application/dto/product"
)

// CreateDiscountCommand represents a command to create a new discount
type CreateDiscountCommand struct {
	Code       *string                   `json:"code" validate:"required,max=100" error:"required=کد تخفیف الزامی است|max=کد تخفیف نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Quantity   *int                      `json:"quantity" validate:"required,max=1000" error:"required=تعداد الزامی است|max=تعداد نمی‌تواند بیشتر از 1000 باشد"`
	Type       *product.DiscountTypeEnum `json:"type" validate:"required" error:"required=نوع تخفیف الزامی است"`
	Value      *int64                    `json:"value" validate:"required" error:"required=مقدار تخفیف الزامی است"`
	ExpiryDate *time.Time                `json:"expiryDate" validate:"required,gtfield=time.Now" error:"required=تاریخ انقضا الزامی است|gtfield=تاریخ انقضا باید در آینده باشد"`
	SiteID     *int64                    `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
}

// UpdateDiscountCommand represents a command to update an existing discount
type UpdateDiscountCommand struct {
	ID         *int64                    `json:"id" validate:"required" error:"required=تخفیف الزامی است"`
	Code       *string                   `json:"code,omitempty" validate:"omitempty,max=100" error:"max=کد تخفیف نمی‌تواند بیشتر از 100 کاراکتر باشد"`
	Quantity   *int                      `json:"quantity,omitempty" validate:"omitempty,gt=0" error:"gt=تعداد باید بزرگتر از 0 باشد"`
	Type       *product.DiscountTypeEnum `json:"type,omitempty" validate:"omitempty" error:""`
	Value      *int64                    `json:"value,omitempty" validate:"omitempty,gt=0" error:"gt=مقدار تخفیف باید بزرگتر از 0 باشد"`
	ExpiryDate *time.Time                `json:"expiryDate,omitempty" validate:"omitempty,gtfield=time.Now" error:"gtfield=تاریخ انقضا باید در آینده باشد"`
}

// DeleteDiscountCommand represents a command to delete a discount
type DeleteDiscountCommand struct {
	ID *int64 `json:"id" validate:"required" error:"required=تخفیف الزامی است"`
}
