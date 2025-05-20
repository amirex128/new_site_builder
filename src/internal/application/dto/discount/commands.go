package discount

import (
	"github.com/amirex128/new_site_builder/src/internal/domain/enums"
	"time"
)

// CreateDiscountCommand represents a command to create a new discount
type CreateDiscountCommand struct {
	Code       *string                 `json:"code" nameFa:"کد" validate:"required_text=1 100"`
	Quantity   *int                    `json:"quantity" nameFa:"تعداد" validate:"required,max=1000"`
	Type       *enums.DiscountTypeEnum `json:"type" nameFa:"نوع" validate:"required,enum"`
	Value      *int64                  `json:"value" nameFa:"مقدار" validate:"required"`
	ExpiryDate *time.Time              `json:"expiryDate" nameFa:"تاریخ انقضا" validate:"required,gtfield=time.Now"`
	SiteID     *int64                  `json:"siteId" nameFa:"شناسه سایت" validate:"required"`
}

// UpdateDiscountCommand represents a command to update an existing discount
type UpdateDiscountCommand struct {
	ID         *int64                  `json:"id" nameFa:"شناسه" validate:"required"`
	Code       *string                 `json:"code,omitempty" nameFa:"کد" validate:"optional_text=1 100"`
	Quantity   *int                    `json:"quantity,omitempty" nameFa:"تعداد" validate:"omitempty,gt=0"`
	Type       *enums.DiscountTypeEnum `json:"type,omitempty" nameFa:"نوع" validate:"enum_optional"`
	Value      *int64                  `json:"value,omitempty" nameFa:"مقدار" validate:"omitempty,gt=0"`
	ExpiryDate *time.Time              `json:"expiryDate,omitempty" nameFa:"تاریخ انقضا" validate:"omitempty,gtfield=time.Now"`
}

// DeleteDiscountCommand represents a command to delete a discount
type DeleteDiscountCommand struct {
	ID *int64 `json:"id" nameFa:"شناسه" validate:"required"`
}
